// Copyright (c) 2015-2022 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio/internal/bucket/lifecycle"
)

//go:generate stringer -type=scannerMetric -trimprefix=scannerMetric $GOFILE

type scannerMetric uint8

type scannerMetrics struct {
	// All fields must be accessed atomically and aligned.
	operations [scannerMetricLast]uint64
	latency    [scannerMetricLastRealtime]lockedLastMinuteLatency

	// actions records actions performed.
	actions        [lifecycle.ActionCount]uint64
	actionsLatency [lifecycle.ActionCount]lockedLastMinuteLatency

	// currentPaths contains (string,*currentPathTracker) for each disk processing.
	// Alignment not required.
	currentPaths sync.Map

	cycleInfoMu sync.Mutex
	cycleInfo   map[string]map[setID]currentScannerCycle
}

var globalScannerMetrics = scannerMetrics{cycleInfo: make(map[string]map[setID]currentScannerCycle)}

const (
	// START Realtime metrics, that only to records
	// last minute latencies and total operation count.
	scannerMetricReadMetadata scannerMetric = iota
	scannerMetricCheckMissing
	scannerMetricSaveUsage
	scannerMetricApplyAll
	scannerMetricApplyVersion
	scannerMetricTierObjSweep
	scannerMetricHealCheck
	scannerMetricILM
	scannerMetricCheckReplication
	scannerMetricYield
	scannerMetricCleanAbandoned
	scannerMetricApplyNonCurrent
	scannerMetricHealAbandonedVersion

	// START Trace metrics:
	scannerMetricStartTrace
	scannerMetricScanObject // Scan object. All operations included.
	scannerMetricHealAbandonedObject

	// END realtime metrics:
	scannerMetricLastRealtime

	// Trace only metrics:
	scannerMetricScanFolder      // Scan a folder on disk, recursively.
	scannerMetricScanCycle       // Full cycle, cluster global
	scannerMetricScanBucketDrive // Single bucket on one drive
	scannerMetricCompactFolder   // Folder compacted.

	// Must be last:
	scannerMetricLast
)

// log scanner action.
// Use for s > scannerMetricStartTrace
func (p *scannerMetrics) log(s scannerMetric, paths ...string) func(custom map[string]string) {
	startTime := time.Now()
	return func(custom map[string]string) {
		duration := time.Since(startTime)

		atomic.AddUint64(&p.operations[s], 1)
		if s < scannerMetricLastRealtime {
			p.latency[s].add(duration)
		}

		if s > scannerMetricStartTrace && globalTrace.NumSubscribers(madmin.TraceScanner) > 0 {
			globalTrace.Publish(scannerTrace(s, startTime, duration, strings.Join(paths, " "), custom))
		}
	}
}

// time a scanner action.
// Use for s < scannerMetricLastRealtime
func (p *scannerMetrics) time(s scannerMetric) func() {
	startTime := time.Now()
	return func() {
		duration := time.Since(startTime)

		atomic.AddUint64(&p.operations[s], 1)
		if s < scannerMetricLastRealtime {
			p.latency[s].add(duration)
		}
	}
}

// timeSize add time and size of a scanner action.
// Use for s < scannerMetricLastRealtime
func (p *scannerMetrics) timeSize(s scannerMetric) func(sz int) {
	startTime := time.Now()
	return func(sz int) {
		duration := time.Since(startTime)

		atomic.AddUint64(&p.operations[s], 1)
		if s < scannerMetricLastRealtime {
			p.latency[s].addSize(duration, int64(sz))
		}
	}
}

// incTime will increment time on metric s with a specific duration.
// Use for s < scannerMetricLastRealtime
func (p *scannerMetrics) incTime(s scannerMetric, d time.Duration) {
	atomic.AddUint64(&p.operations[s], 1)
	if s < scannerMetricLastRealtime {
		p.latency[s].add(d)
	}
}

// timeILM times an ILM action.
// lifecycle.NoneAction is ignored.
// Use for s < scannerMetricLastRealtime
func (p *scannerMetrics) timeILM(a lifecycle.Action) func(versions uint64) {
	if a == lifecycle.NoneAction || a >= lifecycle.ActionCount {
		return func(_ uint64) {}
	}
	startTime := time.Now()
	return func(versions uint64) {
		duration := time.Since(startTime)
		atomic.AddUint64(&p.actions[a], versions)
		p.actionsLatency[a].add(duration)
	}
}

type currentPathTracker struct {
	name *unsafe.Pointer // contains atomically accessed *string
}

// currentPathUpdater provides a lightweight update function for keeping track of
// current objects for each disk.
// Returns a function that can be used to update the current object
// and a function to call to when processing finished.
func (p *scannerMetrics) currentPathUpdater(disk, initial string) (update func(path string), done func()) {
	initialPtr := unsafe.Pointer(&initial)
	tracker := &currentPathTracker{
		name: &initialPtr,
	}

	p.currentPaths.Store(disk, tracker)
	return func(path string) {
			atomic.StorePointer(tracker.name, unsafe.Pointer(&path))
		}, func() {
			p.currentPaths.Delete(disk)
		}
}

// getCurrentPaths returns the paths currently being processed.
func (p *scannerMetrics) getCurrentPaths() []string {
	var res []string
	prefix := globalLocalNodeName + "/"
	p.currentPaths.Range(func(key, value interface{}) bool {
		// We are a bit paranoid, but better miss an entry than crash.
		name, ok := key.(string)
		if !ok {
			return true
		}
		obj, ok := value.(*currentPathTracker)
		if !ok {
			return true
		}
		strptr := (*string)(atomic.LoadPointer(obj.name))
		if strptr != nil {
			res = append(res, pathJoin(prefix, name, *strptr))
		}
		return true
	})
	return res
}

// activeDrives returns the number of currently active disks.
// (since this is concurrent it may not be 100% reliable)
func (p *scannerMetrics) activeDrives() int {
	var i int
	p.currentPaths.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

// lifetime returns the lifetime count of the specified metric.
func (p *scannerMetrics) lifetime(m scannerMetric) uint64 {
	if m >= scannerMetricLast {
		return 0
	}
	val := atomic.LoadUint64(&p.operations[m])
	return val
}

// lastMinute returns the last minute statistics of a metric.
// m should be < scannerMetricLastRealtime
func (p *scannerMetrics) lastMinute(m scannerMetric) AccElem {
	if m >= scannerMetricLastRealtime {
		return AccElem{}
	}
	val := p.latency[m].total()
	return val
}

// lifetimeActions returns the lifetime count of the specified ilm metric.
func (p *scannerMetrics) lifetimeActions(a lifecycle.Action) uint64 {
	if a == lifecycle.NoneAction || a >= lifecycle.ActionCount {
		return 0
	}
	val := atomic.LoadUint64(&p.actions[a])
	return val
}

// lastMinuteActions returns the last minute statistics of an ilm metric.
func (p *scannerMetrics) lastMinuteActions(a lifecycle.Action) AccElem {
	if a == lifecycle.NoneAction || a >= lifecycle.ActionCount {
		return AccElem{}
	}
	val := p.actionsLatency[a].total()
	return val
}

// bucketScanStarted updates the cycle count of the given bucket.
func (p *scannerMetrics) bucketScanStarted(id setID, bucket string, cycle uint32, lastKnownUpdate time.Time) {
	p.cycleInfoMu.Lock()
	defer p.cycleInfoMu.Unlock()

	bucketInfo, ok := p.cycleInfo[bucket]
	if !ok {
		bucketInfo = make(map[setID]currentScannerCycle)
	}

	esInfo := bucketInfo[id]
	esInfo.pool = id.pool
	esInfo.set = id.set
	esInfo.current = uint64(cycle)
	esInfo.lastStarted = time.Now()
	esInfo.lastUpdate = lastKnownUpdate
	esInfo.ongoing = true
	bucketInfo[id] = esInfo

	p.cycleInfo[bucket] = bucketInfo
}

func (p *scannerMetrics) bucketScanFinished(id setID, bucket string) {
	p.cycleInfoMu.Lock()
	defer p.cycleInfoMu.Unlock()

	info := p.cycleInfo[bucket][id]
	info.ongoing = false
	info.cycleCompleted = append(info.cycleCompleted, time.Now())
	if len(info.cycleCompleted) > 10 {
		info.cycleCompleted = info.cycleCompleted[len(info.cycleCompleted)-10:]
	}
	p.cycleInfo[bucket][id] = info
}

func (p *scannerMetrics) bucketScanUpdate(id setID, bucket string) {
	p.cycleInfoMu.Lock()
	defer p.cycleInfoMu.Unlock()

	info := p.cycleInfo[bucket][id]
	info.lastUpdate = time.Now()
	p.cycleInfo[bucket][id] = info
}

type currentScannerCycle struct {
	pool, set      int
	current        uint64
	ongoing        bool
	lastUpdate     time.Time
	lastStarted    time.Time
	cycleCompleted []time.Time // Last 10 completed cycles
}

// clone returns a clone.
func (z currentScannerCycle) clone() currentScannerCycle {
	z.cycleCompleted = append(make([]time.Time, 0, len(z.cycleCompleted)), z.cycleCompleted...)
	return z
}

// getCycles returns the current cycle metrics for all buckets
// If not nil, the returned result can safely be modified.
func (p *scannerMetrics) getCycles() map[string]map[setID]currentScannerCycle {
	p.cycleInfoMu.Lock()
	defer p.cycleInfoMu.Unlock()
	if p.cycleInfo == nil {
		return nil
	}
	ret := make(map[string]map[setID]currentScannerCycle, len(p.cycleInfo))
	for bucket, scanInfo := range p.cycleInfo {
		ret[bucket] = make(map[setID]currentScannerCycle)
		for id, info := range scanInfo {
			ret[bucket][id] = info.clone()
		}
	}
	return ret
}

func (p *scannerMetrics) report() madmin.ScannerMetrics {
	var m madmin.ScannerMetrics

	buckets := make(map[string]struct{})
	for bucket, scanInfo := range p.getCycles() {
		for _, info := range scanInfo {
			if info.ongoing {
				buckets[bucket] = struct{}{}
			}
		}
	}
	m.OngoingBuckets = len(buckets)

	m.CollectedAt = time.Now()
	m.ActivePaths = p.getCurrentPaths()
	m.LifeTimeOps = make(map[string]uint64, scannerMetricLast)
	for i := scannerMetric(0); i < scannerMetricLast; i++ {
		if n := atomic.LoadUint64(&p.operations[i]); n > 0 {
			m.LifeTimeOps[i.String()] = n
		}
	}
	if len(m.LifeTimeOps) == 0 {
		m.LifeTimeOps = nil
	}

	m.LastMinute.Actions = make(map[string]madmin.TimedAction, scannerMetricLastRealtime)
	for i := scannerMetric(0); i < scannerMetricLastRealtime; i++ {
		lm := p.lastMinute(i)
		if lm.N > 0 {
			m.LastMinute.Actions[i.String()] = lm.asTimedAction()
		}
	}
	if len(m.LastMinute.Actions) == 0 {
		m.LastMinute.Actions = nil
	}

	// ILM
	m.LifeTimeILM = make(map[string]uint64)
	for i := lifecycle.NoneAction + 1; i < lifecycle.ActionCount; i++ {
		if n := atomic.LoadUint64(&p.actions[i]); n > 0 {
			m.LifeTimeILM[i.String()] = n
		}
	}
	if len(m.LifeTimeILM) == 0 {
		m.LifeTimeILM = nil
	}

	if len(m.LifeTimeILM) > 0 {
		m.LastMinute.ILM = make(map[string]madmin.TimedAction, len(m.LifeTimeILM))
		for i := lifecycle.NoneAction + 1; i < lifecycle.ActionCount; i++ {
			lm := p.lastMinuteActions(i)
			if lm.N > 0 {
				m.LastMinute.ILM[i.String()] = madmin.TimedAction{Count: uint64(lm.N), AccTime: uint64(lm.Total)}
			}
		}
		if len(m.LastMinute.ILM) == 0 {
			m.LastMinute.ILM = nil
		}
	}
	return m
}
