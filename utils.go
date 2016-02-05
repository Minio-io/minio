/*
 * Minio Cloud Storage, (C) 2015 Minio, Inc.
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

package main

import (
	"encoding/base64"
	"strings"
)

// isValidMD5 - verify if valid md5
func isValidMD5(md5 string) bool {
	if md5 == "" {
		return true
	}
	_, err := base64.StdEncoding.DecodeString(strings.TrimSpace(md5))
	if err != nil {
		return false
	}
	return true
}

/// http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html
const (
	// maximum object size per PUT request is 5GiB
	maxObjectSize = 1024 * 1024 * 1024 * 5
)

// isMaxObjectSize - verify if max object size
func isMaxObjectSize(size int64) bool {
	if size > maxObjectSize {
		return true
	}
	return false
}

func contains(stringList []string, element string) bool {
	for _, e := range stringList {
		if e == element {
			return true
		}
	}

	return false
}
