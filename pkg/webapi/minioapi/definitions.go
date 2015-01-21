/*
 * Mini Object Storage, (C) 2014 Minio, Inc.
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

package minioapi

import (
	"encoding/xml"
)

type ObjectListResponse struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Name        string   `xml:"Name"`
	Marker      string
	MaxKeys     int
	IsTruncated bool
	Contents    []Content `xml:"Contents",innerxml`
}

type BucketListResponse struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   Owner
	Buckets []Bucket `xml:"Buckets",innerxml`
}

type Bucket struct {
	Name         string
	CreationDate string
}

type Content struct {
	Key          string
	LastModified string
	ETag         string
	Size         int
	StorageClass string
	Owner        Owner
}

type Owner struct {
	ID          string
	DisplayName string
}
