/*
 * Mini Object Storage, (C) 2015 Minio, Inc.
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

package singledisk

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"github.com/minio-io/minio/pkg/encoding/erasure"
	"github.com/minio-io/minio/pkg/storage"
	"github.com/minio-io/minio/pkg/storage/donut/erasure/erasure1"
	"github.com/minio-io/minio/pkg/storage/donut/object/objectv1"
	"github.com/minio-io/minio/pkg/utils/split"
	"io"
	"os"
	"path"
	"time"
)

// StorageDriver creates a new single disk storage driver using donut without encoding.
type StorageDriver struct {
	donutBox DonutBox
}

// DonutBox is an interface specifying how the storage driver should interact with its underlying system.
type DonutBox interface {
	Store(objectv1.ObjectMetadata, erasure1.DataHeader, io.Reader) error
	Get(bucket, key string, erasurePart uint16, encodedPart uint8) (objectv1.ObjectMetadata, erasure1.DataHeader, io.Reader, error)
	ListObjects(bucket string) ([]string, error)
	ListBuckets() ([]storage.BucketMetadata, error)
}

// Start a single disk subsystem
func Start(root string, donutBox DonutBox) (chan<- string, <-chan error, storage.Storage) {
	ctrlChannel := make(chan string)
	errorChannel := make(chan error)
	s := new(StorageDriver)
	s.donutBox = donutBox
	go start(ctrlChannel, errorChannel, s)
	return ctrlChannel, errorChannel, s
}

func start(ctrlChannel <-chan string, errorChannel chan<- error, s *StorageDriver) {
	err := os.MkdirAll(s.root, 0700)
	errorChannel <- err
	close(errorChannel)
}

// ListBuckets returns a list of buckets
func (diskStorage StorageDriver) ListBuckets() ([]storage.BucketMetadata, error) {
	return nil, errors.New("Not Implemented")
}

// CreateBucket creates a new bucket
func (diskStorage StorageDriver) CreateBucket(bucket string) error {
	return errors.New("Not Implemented")
}

// CreateBucketPolicy sets a bucket's access policy
func (diskStorage StorageDriver) CreateBucketPolicy(bucket string, p storage.BucketPolicy) error {
	return errors.New("Not Implemented")
}

// GetBucketPolicy returns a bucket's access policy
func (diskStorage StorageDriver) GetBucketPolicy(bucket string) (storage.BucketPolicy, error) {
	return storage.BucketPolicy{}, errors.New("Not Implemented")
}

// GetObject retrieves an object and writes it to a writer
func (diskStorage StorageDriver) GetObject(target io.Writer, bucket, key string) (int64, error) {
	return 0, errors.New("Not Implemented")
}

// GetPartialObject retrieves an object and writes it to a writer
func (diskStorage StorageDriver) GetPartialObject(w io.Writer, bucket, object string, start, length int64) (int64, error) {
	return 0, errors.New("Not Implemented")
}

// GetObjectMetadata retrieves an object's metadata
func (diskStorage StorageDriver) GetObjectMetadata(bucket, key string, prefix string) (metadata storage.ObjectMetadata, err error) {
	return metadata, errors.New("Not Implemented")
}

func readHeaderGob(reader io.Reader) (header ObjectHeader, err error) {
	return header, errors.New("Not Implemented")
}

// ListObjects lists objects
func (diskStorage StorageDriver) ListObjects(bucket string, resources storage.BucketResourcesMetadata) ([]storage.ObjectMetadata, storage.BucketResourcesMetadata, error) {
	return nil, storage.BucketResourcesMetadata{}, errors.New("Not Implemented")
}

// CreateObject creates a new object
func (diskStorage StorageDriver) CreateObject(bucket string, key string, contentType string, data io.Reader) error {
	return errors.New("Not Implemented")
}
