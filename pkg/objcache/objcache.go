/*
 * Minio Cloud Storage, (C) 2016 Minio, Inc.
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
 *
 */

// Package objcache implements in memory caching methods.
package objcache

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"time"
)

// NoExpiry represents caches to be permanent and can only be deleted.
var NoExpiry = time.Duration(0)

// DefaultExpiry represents default time duration value when individual entries will be expired.
var DefaultExpiry = time.Duration(72 * time.Hour) // 72hrs.

// buffer represents the in memory cache of a single entry.
// buffer carries value of the data and last accessed time.
type buffer struct {
	value        []byte    // Value of the entry.
	lastAccessed time.Time // Represents time when value was last accessed.
}

// Cache holds the required variables to compose an in memory cache system
// which also provides expiring key mechanism and also maxSize.
type Cache struct {
	// Mutex is used for handling the concurrent
	// read/write requests for cache
	mutex sync.Mutex

	// maxSize is a total size for overall cache
	maxSize uint64

	// currentSize is a current size in memory
	currentSize uint64

	// OnEviction - callback function for eviction
	OnEviction func(key string)

	// totalEvicted counter to keep track of total expirys
	totalEvicted int

	// map of objectName and its contents
	entries map[string]*buffer

	// Expiry in time duration.
	expiry time.Duration

	// Stop garbage collection routine, stops any running GC routine.
	stopGC chan struct{}
}

// New - Return a new cache with a given default expiry duration.
// If the expiry duration is less than one (or NoExpiry),
// the items in the cache never expire (by default), and must be deleted
// manually.
func New(maxSize uint64, expiry time.Duration) *Cache {
	if maxSize == 0 {
		panic("objcache: setting maximum cache size to zero is forbidden.")
	}
	C := &Cache{
		maxSize: maxSize,
		entries: make(map[string]*buffer),
		expiry:  expiry,
	}
	// We have expiry start the janitor routine.
	if expiry > 0 {
		C.stopGC = make(chan struct{})

		// Start garbage collection routine to expire objects.
		C.startGC()
	}
	return C
}

// ErrKeyNotFoundInCache - key not found in cache.
var ErrKeyNotFoundInCache = errors.New("Key not found in cache")

// ErrCacheFull - cache is full.
var ErrCacheFull = errors.New("Not enough space in cache")

// ErrExcessData - excess data was attempted to be written on cache.
var ErrExcessData = errors.New("Attempted excess write on cache")

// Used for adding entry to the object cache. Implements io.WriteCloser
type cacheBuffer struct {
	*bytes.Buffer // Implements io.Writer
	onClose       func() error
}

// On close, onClose() is called which checks if all object contents
// have been written so that it can save the buffer to the cache.
func (c cacheBuffer) Close() (err error) {
	return c.onClose()
}

// Create - validates if object size fits with in cache size limit and returns a io.WriteCloser
// to which object contents can be written and finally Close()'d. During Close() we
// checks if the amount of data written is equal to the size of the object, in which
// case it saves the contents to object cache.
func (c *Cache) Create(key string, size int64) (w io.WriteCloser, err error) {
	// Recovers any panic generated and return errors appropriately.
	defer func() {
		if r := recover(); r != nil {
			// Recover any panic and return ErrCacheFull.
			err = ErrCacheFull
		}
	}() // Do not crash the server.

	valueLen := uint64(size)
	// Check if the size of the object is not bigger than the capacity of the cache.
	if c.maxSize > 0 && valueLen > c.maxSize {
		return nil, ErrCacheFull
	}

	// Will hold the object contents.
	buf := bytes.NewBuffer(make([]byte, 0, size))

	// Function called on close which saves the object contents
	// to the object cache.
	onClose := func() error {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		if size != int64(buf.Len()) {
			// Full object not available hence do not save buf to object cache.
			return io.ErrShortBuffer
		}
		if c.maxSize > 0 && c.currentSize+valueLen > c.maxSize {
			return ErrExcessData
		}
		// Full object available in buf, save it to cache.
		c.entries[key] = &buffer{
			value:        buf.Bytes(),
			lastAccessed: time.Now().UTC(), // Save last accessed time.
		}
		// Account for the memory allocated above.
		c.currentSize += uint64(size)
		return nil
	}

	// Object contents that is written - cacheBuffer.Write(data)
	// will be accumulated in buf which implements io.Writer.
	return cacheBuffer{
		buf,
		onClose,
	}, nil
}

// Open - open the in-memory file, returns an in memory read seeker.
// returns an error ErrNotFoundInCache, if the key does not exist.
// Returns ErrKeyNotFoundInCache if entry's lastAccessedTime is older
// than objModTime.
func (c *Cache) Open(key string, objModTime time.Time) (io.ReadSeeker, error) {
	// Entry exists, return the readable buffer.
	c.mutex.Lock()
	defer c.mutex.Unlock()
	buf, ok := c.entries[key]
	if !ok {
		return nil, ErrKeyNotFoundInCache
	}
	// Check if buf is recent copy of the object on disk.
	if buf.lastAccessed.Before(objModTime) {
		c.delete(key)
		return nil, ErrKeyNotFoundInCache
	}
	buf.lastAccessed = time.Now().UTC()
	return bytes.NewReader(buf.value), nil
}

// Delete - delete deletes an entry from the cache.
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	c.delete(key)
	c.mutex.Unlock()
	if c.OnEviction != nil {
		c.OnEviction(key)
	}
}

// gc - garbage collect all the expired entries from the cache.
func (c *Cache) gc() {
	var evictedEntries []string
	c.mutex.Lock()
	for k, v := range c.entries {
		if c.expiry > 0 && time.Now().UTC().Sub(v.lastAccessed) > c.expiry {
			c.delete(k)
			evictedEntries = append(evictedEntries, k)
		}
	}
	c.mutex.Unlock()
	for _, k := range evictedEntries {
		if c.OnEviction != nil {
			c.OnEviction(k)
		}
	}
}

// StopGC sends a message to the expiry routine to stop
// expiring cached entries. NOTE: once this is called, cached
// entries will not be expired if the consumer has called this.
func (c *Cache) StopGC() {
	if c.stopGC != nil {
		c.stopGC <- struct{}{}
	}
}

// startGC starts running a routine ticking at expiry interval, on each interval
// this routine does a sweep across the cache entries and garbage collects all the
// expired entries.
func (c *Cache) startGC() {
	go func() {
		for {
			select {
			// Wait till cleanup interval and initiate delete expired entries.
			case <-time.After(c.expiry / 4):
				c.gc()
				// Stop the routine, usually called by the user of object cache during cleanup.
			case <-c.stopGC:
				return
			}
		}
	}()
}

// Deletes a requested entry from the cache.
func (c *Cache) delete(key string) {
	if buf, ok := c.entries[key]; ok {
		delete(c.entries, key)
		c.currentSize -= uint64(len(buf.value))
		c.totalEvicted++
	}
}
