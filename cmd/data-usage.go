/*
 * MinIO Cloud Storage, (C) 2019 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/minio/minio/cmd/config"
	"github.com/minio/minio/cmd/logger"
	"github.com/minio/minio/pkg/color"
	"github.com/minio/minio/pkg/env"
	"github.com/minio/minio/pkg/hash"
)

const (
	dataUsageObjName         = "usage.json"
	dataUsageCacheName       = "usage-cache.bin"
	dataUsageBucketCacheDir  = "usage-caches"
	dataUsageCrawlConf       = "MINIO_DISK_USAGE_CRAWL"
	dataUsageCrawlDelay      = "MINIO_DISK_USAGE_CRAWL_DELAY"
	dataUsageDebug           = true
	dataUsageSleepPerFolder  = 1 * time.Millisecond
	dataUsageSleepDefMult    = 10.0
	dataUsageUpdateDirCycles = 16
	dataUsageRoot            = SlashSeparator
	dataUsageBucket          = bucketConfigPrefix + SlashSeparator + bucketMetaPrefix
)

// initDataUsageStats will start the crawler unless disabled.
func initDataUsageStats() {
	dataUsageEnabled, err := config.ParseBool(env.Get(dataUsageCrawlConf, config.EnableOn))
	if err == nil && !dataUsageEnabled {
		return
	}
	go runDataUsageInfoUpdateRoutine()
}

// runDataUsageInfoUpdateRoutine will contain the main crawler.
func runDataUsageInfoUpdateRoutine() {
	// Wait until the object layer is ready
	var objAPI ObjectLayer
	for {
		objAPI = newObjectLayerWithoutSafeModeFn()
		if objAPI == nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}

	runDataUsageInfo(GlobalContext, objAPI)
}

var dataUsageLockTimeout = lifecycleLockTimeout

func runDataUsageInfo(ctx context.Context, objAPI ObjectLayer) {
	// Make sure only 1 crawler is running on the cluster.
	locker := objAPI.NewNSLock(ctx, minioMetaBucket, "leader-data-usage-info")
	for {
		err := locker.GetLock(dataUsageLockTimeout)
		if err != nil {
			time.Sleep(5 * time.Minute)
			continue
		}
		// Break without unlocking, this node will acquire
		// data usage calculator role for its lifetime.
		break
	}
	if dataUsageDebug {
		logger.Info(color.Green("runDataUsageInfo:") + " Starting crawler master")
	}
	for {
		select {
		case <-ctx.Done():
			locker.Unlock()
			return
			// Wait 5 minutes before starting.
			// FIXME: For testing, set to 5 seconds...
		case <-time.NewTimer(5 * time.Second).C:
			results := make(chan DataUsageInfo, 1)
			go storeDataUsageInBackend(ctx, objAPI, results)
			err := objAPI.CrawlAndGetDataUsage(ctx, results)
			close(results)
			logger.LogIf(ctx, err)
		}
	}
}

// storeDataUsageInBackend will store all objects sent on the gui channel until closed.
func storeDataUsageInBackend(ctx context.Context, objAPI ObjectLayer, gui <-chan DataUsageInfo) {
	for dataUsageInfo := range gui {
		dataUsageJSON, err := json.MarshalIndent(dataUsageInfo, "", "  ")
		if err != nil {
			logger.LogIf(ctx, err)
			continue
		}
		if dataUsageDebug {
			logger.Info(color.Green("data-usage:")+" Received update: %s", string(dataUsageJSON))
		}
		size := int64(len(dataUsageJSON))
		r, err := hash.NewReader(bytes.NewReader(dataUsageJSON), size, "", "", size, false)
		if err != nil {
			logger.LogIf(ctx, err)
			continue
		}

		_, err = objAPI.PutObject(ctx, dataUsageBucket, dataUsageObjName, NewPutObjReader(r, nil, nil), ObjectOptions{})
		logger.LogIf(ctx, err)
	}
}

func loadDataUsageFromBackend(ctx context.Context, objAPI ObjectLayer) (DataUsageInfo, error) {
	var dataUsageInfoJSON bytes.Buffer

	err := objAPI.GetObject(ctx, dataUsageBucket, dataUsageObjName, 0, -1, &dataUsageInfoJSON, "", ObjectOptions{})
	if err != nil {
		if isErrObjectNotFound(err) {
			return DataUsageInfo{}, nil
		}
		return DataUsageInfo{}, toObjectErr(err, dataUsageBucket, dataUsageObjName)
	}

	var dataUsageInfo DataUsageInfo
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(dataUsageInfoJSON.Bytes(), &dataUsageInfo)
	if err != nil {
		return DataUsageInfo{}, err
	}

	return dataUsageInfo, nil
}

// Item represents each file while walking.
type Item struct {
	Path string
	Typ  os.FileMode
}

type getSizeFn func(item Item) (int64, error)

type cachedFolder struct {
	name   string
	parent *dataUsageHash
}

type folderScanner struct {
	root               string
	getSize            getSizeFn
	oldCache           dataUsageCache
	newCache           dataUsageCache
	waitForLowActiveIO func()

	newFolders      []cachedFolder
	existingFolders []cachedFolder
}

// sleepDuration multiplies the duration d by x and sleeps if is more than 100 micro seconds.
// sleep is limited to max 1 second.
func sleepDuration(d time.Duration, x float64) {
	// Don't sleep for really small amount of time
	if d := time.Duration(float64(d) * x); d > time.Microsecond*100 {
		if d > time.Second {
			d = time.Second
		}
		time.Sleep(d)
	}
}

// scanQueuedLevels will scan the provided folders.
// Files found in the folders will be added to f.newCache.
// If final is provided folders will be put into f.newFolders or f.existingFolders.
// If final is not provided the folders found are returned from the function.
func (f *folderScanner) scanQueuedLevels(ctx context.Context, folders []cachedFolder, final bool) ([]cachedFolder, error) {
	var nextFolders []cachedFolder
	delayMult := dataUsageSleepDefMult
	if mult := os.Getenv(dataUsageCrawlDelay); mult != "" {
		if d, err := strconv.ParseFloat(mult, 64); err == nil {
			delayMult = d
		}
	}
	done := ctx.Done()
	for _, folder := range folders {
		select {
		case <-done:
			return nil, ctx.Err()
		default:
		}
		f.waitForLowActiveIO()
		sleepDuration(dataUsageSleepPerFolder, delayMult)

		cache := dataUsageEntry{}
		thisHash := hashPath(folder.name)

		err := readDirFn(path.Join(f.root, folder.name), func(entName string, typ os.FileMode) error {
			// Parse
			entName = path.Clean(path.Join(folder.name, entName))
			bucket, _ := path2BucketObjectWithBasePath(f.root, entName)
			if bucket == "" {
				if dataUsageDebug {
					logger.Info(color.Green("data-usage:")+" no bucket (%s,%s)", f.root, entName)
				}
				return nil
			}

			if isReservedOrInvalidBucket(bucket, false) {
				if dataUsageDebug {
					logger.Info(color.Green("data-usage:")+" invalid bucket: %v, entry: %v", bucket, entName)
				}
				return nil
			}

			select {
			case <-done:
				return ctx.Err()
			default:
			}

			if typ&os.ModeDir != 0 {
				h := hashPath(entName)
				_, exists := f.oldCache.Cache[h]
				cache.addChildString(entName)

				this := cachedFolder{name: entName, parent: &thisHash}
				cache.addChild(h)
				if final {
					if exists {
						f.existingFolders = append(f.existingFolders, this)
					} else {
						f.newFolders = append(f.newFolders, this)
					}
				} else {
					nextFolders = append(nextFolders, this)
				}
				return nil
			}
			f.waitForLowActiveIO()
			// Dynamic time delay.
			t := time.Now()

			// Get file size, ignore errors.
			size, err := f.getSize(Item{Path: path.Join(f.root, entName), Typ: typ})

			sleepDuration(time.Since(t), delayMult)
			if err == errSkipFile {
				return nil
			}
			logger.LogIf(ctx, err)
			cache.Size += size
			cache.Objects++
			cache.ObjSizes.add(size)

			return nil
		})
		if err != nil {
			return nil, err
		}
		f.newCache.replaceHashed(thisHash, folder.parent, cache)
	}
	return nextFolders, nil
}

// deepScanFolder will deep scan a folder and return the size if no error occurs.
func (f *folderScanner) deepScanFolder(ctx context.Context, folder string) (*dataUsageEntry, error) {
	var cache dataUsageEntry
	delayMult := dataUsageSleepDefMult
	if mult := os.Getenv(dataUsageCrawlDelay); mult != "" {
		if d, err := strconv.ParseFloat(mult, 64); err == nil {
			delayMult = d
		}
	}
	done := ctx.Done()

	var addDir func(entName string, typ os.FileMode) error
	var dirStack = []string{f.root, folder}

	addDir = func(entName string, typ os.FileMode) error {
		select {
		case <-done:
			return ctx.Err()
		default:
		}

		f.waitForLowActiveIO()
		if typ&os.ModeDir != 0 {
			dirStack = append(dirStack, entName)
			err := readDirFn(path.Join(dirStack...), addDir)
			dirStack = dirStack[:len(dirStack)-1]
			sleepDuration(dataUsageSleepPerFolder, delayMult)
			return err
		}
		// Dynamic time delay.
		t := time.Now()

		// Get file size, ignore errors.
		dirStack = append(dirStack, entName)
		fileName := path.Join(dirStack...)
		dirStack = dirStack[:len(dirStack)-1]

		size, err := f.getSize(Item{Path: fileName, Typ: typ})

		// Don't sleep for really small amount of time
		sleepDuration(time.Since(t), delayMult)

		if err == errSkipFile {
			return nil
		}
		logger.LogIf(ctx, err)
		cache.Size += size
		cache.Objects++
		cache.ObjSizes.add(size)
		return nil
	}
	err := readDirFn(path.Join(dirStack...), addDir)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// updateUsage will crawl the basepath+cache.Info.Name and return an updated cache.
// The returned cache will always be valid, but may not be updated from the existing.
// Before each operation waitForLowActiveIO is called which can be used to temporarily halt the crawler.
// If the supplied context is canceled the function will return at the first chance.
func updateUsage(ctx context.Context, basePath string, cache dataUsageCache, waitForLowActiveIO func(), getSize getSizeFn) (dataUsageCache, error) {
	if cache.Info.Name == "" {
		cache.Info.Name = dataUsageRoot
	}
	var logPrefix, logSuffix string
	if dataUsageDebug {
		logPrefix = color.Green("data-usage: ")
		logSuffix = color.Blue(" - %v + %v", basePath, cache.Info.Name)
	}
	s := folderScanner{
		root:               basePath,
		getSize:            getSize,
		oldCache:           cache,
		newCache:           dataUsageCache{Info: cache.Info},
		waitForLowActiveIO: waitForLowActiveIO,
		newFolders:         nil,
		existingFolders:    nil,
	}

	done := ctx.Done()
	var flattenLevels = 3

	// If we are scanning inside a bucket reduce depth by 1.
	if cache.Info.Name != dataUsageRoot {
		flattenLevels--
	}
	if dataUsageDebug {
		logger.Info(logPrefix+"Cycle: %v"+logSuffix, cache.Info.NextCycle)
	}

	// Always scan flattenLevels deep. Cache root is level 0.
	todo := []cachedFolder{{name: cache.Info.Name}}
	for i := 0; i < flattenLevels; i++ {
		if dataUsageDebug {
			logger.Info(logPrefix+"Level %v, scanning %v directories."+logSuffix, i, len(todo))
		}
		select {
		case <-done:
			return cache, ctx.Err()
		default:
		}
		var err error
		todo, err = s.scanQueuedLevels(ctx, todo, i == flattenLevels-1)
		if err != nil {
			// No useful information...
			return cache, err
		}
	}

	if dataUsageDebug {
		logger.Info(logPrefix+"New folders: %v"+logSuffix, s.newFolders)
	}
	// Add new folders first
	for _, folder := range s.newFolders {
		select {
		case <-done:
			return s.newCache, ctx.Err()
		default:
		}
		du, err := s.deepScanFolder(ctx, folder.name)
		if err != nil {
			logger.LogIf(ctx, err)
			continue
		}
		if du == nil {
			logger.LogIf(ctx, errors.New("data-usage: no disk usage provided"))
			continue
		}
		s.newCache.replace(folder.name, "", *du)
		// Add to parent manually
		if folder.parent != nil {
			parent := s.newCache.Cache[*folder.parent]
			parent.addChildString(folder.name)
		}
	}

	if dataUsageDebug {
		logger.Info(logPrefix+"Existing folders: %v"+logSuffix, len(s.existingFolders))
	}
	// Do selective scanning of existing folders.
	for _, folder := range s.existingFolders {
		select {
		case <-done:
			return s.newCache, ctx.Err()
		default:
		}
		h := hashPath(folder.name)
		if !h.mod(s.oldCache.Info.NextCycle, dataUsageUpdateDirCycles) {
			s.newCache.replaceHashed(h, folder.parent, s.oldCache.Cache[h])
			continue
		}

		// Update on this cycle...
		du, err := s.deepScanFolder(ctx, folder.name)
		if err != nil {
			logger.LogIf(ctx, err)
			continue
		}
		if du == nil {
			logger.LogIf(ctx, errors.New("data-usage: no disk usage provided"))
			continue
		}
		s.newCache.replaceHashed(h, folder.parent, *du)
	}

	s.newCache.Info.LastUpdate = time.Now()
	s.newCache.Info.NextCycle++
	return s.newCache, nil
}
