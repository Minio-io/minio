package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/minio/madmin-go/v2"
)

const bootstrapMsgsLimit = 4 << 10

type bootstrapInfo struct {
	msg    string
	ts     time.Time
	source string
}
type bootstrapTracer struct {
	mu         sync.RWMutex
	idx        int
	info       [bootstrapMsgsLimit]bootstrapInfo
	lastUpdate time.Time
}

var globalBootstrapTracer = &bootstrapTracer{}

func (bs *bootstrapTracer) DropEvents() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if time.Now().UTC().Sub(bs.lastUpdate) > 24*time.Hour {
		bs.info = [4096]bootstrapInfo{}
		bs.idx = 0
	}
}

func (bs *bootstrapTracer) Empty() bool {
	var empty bool
	bs.mu.RLock()
	empty = bs.info[0].msg == ""
	bs.mu.RUnlock()

	return empty
}

func (bs *bootstrapTracer) Record(msg string) {
	source := getSource(2)
	bs.mu.Lock()
	now := time.Now().UTC()
	bs.info[bs.idx] = bootstrapInfo{
		msg:    msg,
		ts:     now,
		source: source,
	}
	bs.lastUpdate = now
	bs.idx++
	if bs.idx >= bootstrapMsgsLimit {
		bs.idx = 0 // circular buffer
	}
	bs.mu.Unlock()
}

func (bs *bootstrapTracer) Events() []madmin.TraceInfo {
	var info [bootstrapMsgsLimit]bootstrapInfo
	var idx int

	bs.mu.RLock()
	idx = bs.idx
	tail := bootstrapMsgsLimit - idx
	copy(info[tail:], bs.info[:idx])
	copy(info[:tail], bs.info[idx:])
	bs.mu.RUnlock()

	traceInfo := make([]madmin.TraceInfo, 0, bootstrapMsgsLimit)
	for i := 0; i < bootstrapMsgsLimit; i++ {
		if info[i].ts.IsZero() {
			continue // skip empty events
		}
		traceInfo = append(traceInfo, madmin.TraceInfo{
			TraceType: madmin.TraceBootstrap,
			Time:      info[i].ts,
			NodeName:  globalLocalNodeName,
			FuncName:  "BOOTSTRAP",
			Message:   fmt.Sprintf("%s %s", info[i].source, info[i].msg),
		})
	}
	return traceInfo
}
