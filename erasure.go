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
 */

package main

import (
	"errors"

	"github.com/klauspost/reedsolomon"
)

// erasure storage layer.
type erasure struct {
	ReedSolomon  reedsolomon.Encoder // Erasure encoder/decoder.
	DataBlocks   int
	ParityBlocks int
	storageDisks []StorageAPI
}

// errUnexpected - returned for any unexpected error.
var errUnexpected = errors.New("Unexpected error - please report at https://github.com/minio/minio/issues")

// newErasure instantiate a new erasure.
func newErasure(disks []StorageAPI) (*erasure, error) {
	// Initialize E.
	e := &erasure{}

	// Calculate data and parity blocks.
	dataBlocks, parityBlocks := len(disks)/2, len(disks)/2

	// Initialize reed solomon encoding.
	rs, err := reedsolomon.New(dataBlocks, parityBlocks)
	if err != nil {
		return nil, err
	}

	// Save the reedsolomon.
	e.DataBlocks = dataBlocks
	e.ParityBlocks = parityBlocks
	e.ReedSolomon = rs

	// Save all the initialized storage disks.
	e.storageDisks = disks

	// Return successfully initialized.
	return e, nil
}
