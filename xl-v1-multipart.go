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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio/pkg/mimedb"
	"github.com/skyrings/skyring-common/tools/uuid"
)

// listMultipartUploads - lists all multipart uploads.
func (xl xlObjects) listMultipartUploads(bucket, prefix, keyMarker, uploadIDMarker, delimiter string, maxUploads int) (ListMultipartsInfo, error) {
	result := ListMultipartsInfo{
		IsTruncated: true,
		MaxUploads:  maxUploads,
		KeyMarker:   keyMarker,
		Prefix:      prefix,
		Delimiter:   delimiter,
	}

	recursive := true
	if delimiter == slashSeparator {
		recursive = false
	}

	// Not using path.Join() as it strips off the trailing '/'.
	multipartPrefixPath := pathJoin(mpartMetaPrefix, bucket, prefix)
	if prefix == "" {
		// Should have a trailing "/" if prefix is ""
		// For ex. multipartPrefixPath should be "multipart/bucket/" if prefix is ""
		multipartPrefixPath += slashSeparator
	}
	multipartMarkerPath := ""
	if keyMarker != "" {
		multipartMarkerPath = pathJoin(mpartMetaPrefix, bucket, keyMarker)
	}
	var uploads []uploadMetadata
	var err error
	var eof bool
	// List all upload ids for the keyMarker starting from
	// uploadIDMarker first.
	if uploadIDMarker != "" {
		nsMutex.RLock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, keyMarker))
		for _, disk := range xl.getLoadBalancedQuorumDisks() {
			if disk == nil {
				continue
			}
			uploads, _, err = listMultipartUploadIDs(bucket, keyMarker, uploadIDMarker, maxUploads, disk)
			if err == errDiskNotFound {
				continue
			}
			break
		}
		nsMutex.RUnlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, keyMarker))
		if err != nil {
			return ListMultipartsInfo{}, err
		}
		maxUploads = maxUploads - len(uploads)
	}
	// Validate if we need to list further depending on maxUploads.
	if maxUploads > 0 {
		walkerCh, walkerDoneCh := xl.listPool.Release(listParams{minioMetaBucket, recursive, multipartMarkerPath, multipartPrefixPath})
		if walkerCh == nil {
			walkerDoneCh = make(chan struct{})
			walkerCh = xl.startTreeWalk(minioMetaBucket, multipartPrefixPath, multipartMarkerPath, recursive, xl.isMultipartUpload, walkerDoneCh)
		}
		// Collect uploads until we have reached maxUploads count to 0.
		for maxUploads > 0 {
			walkResult, ok := <-walkerCh
			if !ok {
				// Closed channel.
				eof = true
				break
			}
			// For any walk error return right away.
			if walkResult.err != nil {
				// File not found or Disk not found is a valid case.
				if walkResult.err == errFileNotFound || walkResult.err == errDiskNotFound {
					continue
				}
				return ListMultipartsInfo{}, err
			}
			entry := strings.TrimPrefix(walkResult.entry, retainSlash(pathJoin(mpartMetaPrefix, bucket)))
			// For an entry looking like a directory, store and
			// continue the loop not need to fetch uploads.
			if strings.HasSuffix(walkResult.entry, slashSeparator) {
				uploads = append(uploads, uploadMetadata{
					Object: entry,
				})
				maxUploads--
				if maxUploads == 0 {
					eof = true
					break
				}
				continue
			}
			var newUploads []uploadMetadata
			var end bool
			uploadIDMarker = ""
			// For the new object entry we get all its pending uploadIDs.
			nsMutex.RLock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, entry))
			var disk StorageAPI
			for _, disk = range xl.getLoadBalancedQuorumDisks() {
				if disk == nil {
					continue
				}
				newUploads, end, err = listMultipartUploadIDs(bucket, entry, uploadIDMarker, maxUploads, disk)
				if err == errDiskNotFound {
					continue
				}
				break
			}
			nsMutex.RUnlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, entry))
			if err != nil {
				if err == errFileNotFound || walkResult.err == errDiskNotFound {
					continue
				}
				return ListMultipartsInfo{}, err
			}
			uploads = append(uploads, newUploads...)
			maxUploads -= len(newUploads)
			if end && walkResult.end {
				eof = true
				break
			}
		}
	}
	// For all received uploads fill in the multiparts result.
	for _, upload := range uploads {
		var objectName string
		var uploadID string
		if strings.HasSuffix(upload.Object, slashSeparator) {
			// All directory entries are common prefixes.
			uploadID = "" // For common prefixes, upload ids are empty.
			objectName = upload.Object
			result.CommonPrefixes = append(result.CommonPrefixes, objectName)
		} else {
			uploadID = upload.UploadID
			objectName = upload.Object
			result.Uploads = append(result.Uploads, upload)
		}
		result.NextKeyMarker = objectName
		result.NextUploadIDMarker = uploadID
	}
	result.IsTruncated = !eof
	// Result is not truncated, reset the markers.
	if !result.IsTruncated {
		result.NextKeyMarker = ""
		result.NextUploadIDMarker = ""
	}
	return result, nil
}

// ListMultipartUploads - lists all the pending multipart uploads on a
// bucket. Additionally takes 'prefix, keyMarker, uploadIDmarker and a
// delimiter' which allows us to list uploads match a particular
// prefix or lexically starting from 'keyMarker' or delimiting the
// output to get a directory like listing.
//
// Implements S3 compatible ListMultipartUploads API. The resulting
// ListMultipartsInfo structure is unmarshalled directly into XML and
// replied back to the client.
func (xl xlObjects) ListMultipartUploads(bucket, prefix, keyMarker, uploadIDMarker, delimiter string, maxUploads int) (ListMultipartsInfo, error) {
	result := ListMultipartsInfo{}

	// Verify if bucket is valid.
	if !IsValidBucketName(bucket) {
		return ListMultipartsInfo{}, BucketNameInvalid{Bucket: bucket}
	}
	if !xl.isBucketExist(bucket) {
		return ListMultipartsInfo{}, BucketNotFound{Bucket: bucket}
	}
	if !IsValidObjectPrefix(prefix) {
		return ListMultipartsInfo{}, ObjectNameInvalid{Bucket: bucket, Object: prefix}
	}
	// Verify if delimiter is anything other than '/', which we do not support.
	if delimiter != "" && delimiter != slashSeparator {
		return ListMultipartsInfo{}, UnsupportedDelimiter{
			Delimiter: delimiter,
		}
	}
	// Verify if marker has prefix.
	if keyMarker != "" && !strings.HasPrefix(keyMarker, prefix) {
		return ListMultipartsInfo{}, InvalidMarkerPrefixCombination{
			Marker: keyMarker,
			Prefix: prefix,
		}
	}
	if uploadIDMarker != "" {
		if strings.HasSuffix(keyMarker, slashSeparator) {
			return result, InvalidUploadIDKeyCombination{
				UploadIDMarker: uploadIDMarker,
				KeyMarker:      keyMarker,
			}
		}
		id, err := uuid.Parse(uploadIDMarker)
		if err != nil {
			return result, err
		}
		if id.IsZero() {
			return result, MalformedUploadID{
				UploadID: uploadIDMarker,
			}
		}
	}
	return xl.listMultipartUploads(bucket, prefix, keyMarker, uploadIDMarker, delimiter, maxUploads)
}

// newMultipartUpload - wrapper for initializing a new multipart
// request, returns back a unique upload id.
//
// Internally this function creates 'uploads.json' associated for the
// incoming object at '.minio/multipart/bucket/object/uploads.json' on
// all the disks. `uploads.json` carries metadata regarding on going
// multipart operation on the object.
func (xl xlObjects) newMultipartUpload(bucket string, object string, meta map[string]string) (uploadID string, err error) {
	xlMeta := newXLMetaV1(xl.dataBlocks, xl.parityBlocks)
	// If not set default to "application/octet-stream"
	if meta["content-type"] == "" {
		contentType := "application/octet-stream"
		if objectExt := filepath.Ext(object); objectExt != "" {
			content, ok := mimedb.DB[strings.ToLower(strings.TrimPrefix(objectExt, "."))]
			if ok {
				contentType = content.ContentType
			}
		}
		meta["content-type"] = contentType
	}
	xlMeta.Stat.ModTime = time.Now().UTC()
	xlMeta.Stat.Version = 1
	xlMeta.Meta = meta

	// This lock needs to be held for any changes to the directory contents of ".minio/multipart/object/"
	nsMutex.Lock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object))
	defer nsMutex.Unlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object))

	uploadID = getUUID()
	initiated := time.Now().UTC()
	// Create 'uploads.json'
	if err = writeUploadJSON(bucket, object, uploadID, initiated, xl.storageDisks...); err != nil {
		return "", err
	}
	uploadIDPath := path.Join(mpartMetaPrefix, bucket, object, uploadID)
	tempUploadIDPath := path.Join(tmpMetaPrefix, uploadID)
	// Write updated `xl.json` to all disks.
	if err = xl.writeSameXLMetadata(minioMetaBucket, tempUploadIDPath, xlMeta); err != nil {
		return "", toObjectErr(err, minioMetaBucket, tempUploadIDPath)
	}
	rErr := xl.renameObject(minioMetaBucket, tempUploadIDPath, minioMetaBucket, uploadIDPath)
	if rErr == nil {
		// Return success.
		return uploadID, nil
	}
	return "", toObjectErr(rErr, minioMetaBucket, uploadIDPath)
}

// NewMultipartUpload - initialize a new multipart upload, returns a
// unique id. The unique id returned here is of UUID form, for each
// subsequent request each UUID is unique.
//
// Implements S3 compatible initiate multipart API.
func (xl xlObjects) NewMultipartUpload(bucket, object string, meta map[string]string) (string, error) {
	// Verify if bucket name is valid.
	if !IsValidBucketName(bucket) {
		return "", BucketNameInvalid{Bucket: bucket}
	}
	// Verify whether the bucket exists.
	if !xl.isBucketExist(bucket) {
		return "", BucketNotFound{Bucket: bucket}
	}
	// Verify if object name is valid.
	if !IsValidObjectName(object) {
		return "", ObjectNameInvalid{Bucket: bucket, Object: object}
	}
	// No metadata is set, allocate a new one.
	if meta == nil {
		meta = make(map[string]string)
	}
	return xl.newMultipartUpload(bucket, object, meta)
}

// putObjectPart - reads incoming data until EOF for the part file on
// an ongoing multipart transaction. Internally incoming data is
// erasure coded and written across all disks.
func (xl xlObjects) putObjectPart(bucket string, object string, uploadID string, partID int, size int64, data io.Reader, md5Hex string) (string, error) {
	// Hold the lock and start the operation.
	uploadIDPath := pathJoin(mpartMetaPrefix, bucket, object, uploadID)
	nsMutex.Lock(minioMetaBucket, uploadIDPath)
	defer nsMutex.Unlock(minioMetaBucket, uploadIDPath)

	if !xl.isUploadIDExists(bucket, object, uploadID) {
		return "", InvalidUploadID{UploadID: uploadID}
	}

	// Read metadata associated with the object from all disks.
	partsMetadata, errs := xl.readAllXLMetadata(minioMetaBucket, uploadIDPath)

	// List all online disks.
	onlineDisks, higherVersion, err := xl.listOnlineDisks(partsMetadata, errs)
	if err != nil {
		return "", toObjectErr(err, bucket, object)
	}

	// Pick one from the first valid metadata.
	xlMeta := pickValidXLMeta(partsMetadata)

	partSuffix := fmt.Sprintf("object%d", partID)
	tmpPartPath := path.Join(tmpMetaPrefix, uploadID, partSuffix)

	// Initialize md5 writer.
	md5Writer := md5.New()

	// Construct a tee reader for md5sum.
	teeReader := io.TeeReader(data, md5Writer)

	// Collect all the previous erasure infos across the disk.
	var eInfos []erasureInfo
	for index := range onlineDisks {
		eInfos = append(eInfos, partsMetadata[index].Erasure)
	}

	// Erasure code data and write across all disks.
	newEInfos, n, err := erasureCreateFile(onlineDisks, minioMetaBucket, tmpPartPath, partSuffix, teeReader, eInfos)
	if err != nil {
		return "", toObjectErr(err, minioMetaBucket, tmpPartPath)
	}
	// Calculate new md5sum.
	newMD5Hex := hex.EncodeToString(md5Writer.Sum(nil))
	if md5Hex != "" {
		if newMD5Hex != md5Hex {
			// MD5 mismatch, delete the temporary object.
			xl.deleteObject(minioMetaBucket, tmpPartPath)
			// Returns md5 mismatch.
			return "", BadDigest{md5Hex, newMD5Hex}
		}
	}
	if size != n {
		// content-length mismatch, delete the temporary object.
		xl.deleteObject(minioMetaBucket, tmpPartPath)
		return "", errSignatureMismatch
	}

	// Validates if upload ID exists again.
	if !xl.isUploadIDExists(bucket, object, uploadID) {
		return "", InvalidUploadID{UploadID: uploadID}
	}

	// Rename temporary part file to its final location.
	partPath := path.Join(uploadIDPath, partSuffix)
	err = xl.renameObject(minioMetaBucket, tmpPartPath, minioMetaBucket, partPath)
	if err != nil {
		return "", toObjectErr(err, minioMetaBucket, partPath)
	}

	// Once part is successfully committed, proceed with updating XL metadata.
	xlMeta.Stat.Version = higherVersion
	// Add the current part.
	xlMeta.AddObjectPart(partID, partSuffix, newMD5Hex, size)

	// Update `xl.json` content for each disks.
	for index := range partsMetadata {
		partsMetadata[index].Parts = xlMeta.Parts
		partsMetadata[index].Erasure = newEInfos[index]
	}

	// Write all the checksum metadata.
	tempUploadIDPath := path.Join(tmpMetaPrefix, uploadID)

	// Writes a unique `xl.json` each disk carrying new checksum
	// related information.
	if err = xl.writeUniqueXLMetadata(minioMetaBucket, tempUploadIDPath, partsMetadata); err != nil {
		return "", toObjectErr(err, minioMetaBucket, tempUploadIDPath)
	}
	rErr := xl.renameXLMetadata(minioMetaBucket, tempUploadIDPath, minioMetaBucket, uploadIDPath)
	if rErr != nil {
		return "", toObjectErr(rErr, minioMetaBucket, uploadIDPath)
	}

	// Return success.
	return newMD5Hex, nil
}

// PutObjectPart - reads incoming stream and internally erasure codes
// them. This call is similar to single put operation but it is part
// of the multipart transcation.
//
// Implements S3 compatible Upload Part API.
func (xl xlObjects) PutObjectPart(bucket, object, uploadID string, partID int, size int64, data io.Reader, md5Hex string) (string, error) {
	// Verify if bucket is valid.
	if !IsValidBucketName(bucket) {
		return "", BucketNameInvalid{Bucket: bucket}
	}
	// Verify whether the bucket exists.
	if !xl.isBucketExist(bucket) {
		return "", BucketNotFound{Bucket: bucket}
	}
	if !IsValidObjectName(object) {
		return "", ObjectNameInvalid{Bucket: bucket, Object: object}
	}
	return xl.putObjectPart(bucket, object, uploadID, partID, size, data, md5Hex)
}

// listObjectParts - wrapper reading `xl.json` for a given object and
// uploadID. Lists all the parts captured inside `xl.json` content.
func (xl xlObjects) listObjectParts(bucket, object, uploadID string, partNumberMarker, maxParts int) (ListPartsInfo, error) {
	result := ListPartsInfo{}

	uploadIDPath := path.Join(mpartMetaPrefix, bucket, object, uploadID)

	xlMeta, err := xl.readXLMetadata(minioMetaBucket, uploadIDPath)
	if err != nil {
		return ListPartsInfo{}, toObjectErr(err, minioMetaBucket, uploadIDPath)
	}

	// Populate the result stub.
	result.Bucket = bucket
	result.Object = object
	result.UploadID = uploadID
	result.MaxParts = maxParts

	// For empty number of parts or maxParts as zero, return right here.
	if len(xlMeta.Parts) == 0 || maxParts == 0 {
		return result, nil
	}

	// Limit output to maxPartsList.
	if maxParts > maxPartsList {
		maxParts = maxPartsList
	}

	// Only parts with higher part numbers will be listed.
	partIdx := xlMeta.ObjectPartIndex(partNumberMarker)
	parts := xlMeta.Parts
	if partIdx != -1 {
		parts = xlMeta.Parts[partIdx+1:]
	}
	count := maxParts
	for _, part := range parts {
		var fi FileInfo
		fi, err = xl.statPart(bucket, object, uploadID, part.Name)
		if err != nil {
			return ListPartsInfo{}, toObjectErr(err, minioMetaBucket, path.Join(uploadID, part.Name))
		}
		result.Parts = append(result.Parts, partInfo{
			PartNumber:   part.Number,
			ETag:         part.ETag,
			LastModified: fi.ModTime,
			Size:         part.Size,
		})
		count--
		if count == 0 {
			break
		}
	}
	// If listed entries are more than maxParts, we set IsTruncated as true.
	if len(parts) > len(result.Parts) {
		result.IsTruncated = true
		// Make sure to fill next part number marker if IsTruncated is
		// true for subsequent listing.
		nextPartNumberMarker := result.Parts[len(result.Parts)-1].PartNumber
		result.NextPartNumberMarker = nextPartNumberMarker
	}
	return result, nil
}

// ListObjectParts - lists all previously uploaded parts for a given
// object and uploadID.  Takes additional input of part-number-marker
// to indicate where the listing should begin from.
//
// Implements S3 compatible ListObjectParts API. The resulting
// ListPartsInfo structure is unmarshalled directly into XML and
// replied back to the client.
func (xl xlObjects) ListObjectParts(bucket, object, uploadID string, partNumberMarker, maxParts int) (ListPartsInfo, error) {
	// Verify if bucket is valid.
	if !IsValidBucketName(bucket) {
		return ListPartsInfo{}, BucketNameInvalid{Bucket: bucket}
	}
	// Verify whether the bucket exists.
	if !xl.isBucketExist(bucket) {
		return ListPartsInfo{}, BucketNotFound{Bucket: bucket}
	}
	if !IsValidObjectName(object) {
		return ListPartsInfo{}, ObjectNameInvalid{Bucket: bucket, Object: object}
	}
	// Hold lock so that there is no competing abort-multipart-upload or complete-multipart-upload.
	nsMutex.Lock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object, uploadID))
	defer nsMutex.Unlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object, uploadID))

	if !xl.isUploadIDExists(bucket, object, uploadID) {
		return ListPartsInfo{}, InvalidUploadID{UploadID: uploadID}
	}
	result, err := xl.listObjectParts(bucket, object, uploadID, partNumberMarker, maxParts)
	return result, err
}

// CompleteMultipartUpload - completes an ongoing multipart
// transaction after receiving all the parts indicated by the client.
// Returns an md5sum calculated by concatenating all the individual
// md5sums of all the parts.
//
// Implements S3 compatible Complete multipart API.
func (xl xlObjects) CompleteMultipartUpload(bucket string, object string, uploadID string, parts []completePart) (string, error) {
	// Verify if bucket is valid.
	if !IsValidBucketName(bucket) {
		return "", BucketNameInvalid{Bucket: bucket}
	}
	// Verify whether the bucket exists.
	if !xl.isBucketExist(bucket) {
		return "", BucketNotFound{Bucket: bucket}
	}
	if !IsValidObjectName(object) {
		return "", ObjectNameInvalid{
			Bucket: bucket,
			Object: object,
		}
	}
	// Hold lock so that
	// 1) no one aborts this multipart upload
	// 2) no one does a parallel complete-multipart-upload on this multipart upload
	nsMutex.Lock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object, uploadID))
	defer nsMutex.Unlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object, uploadID))

	if !xl.isUploadIDExists(bucket, object, uploadID) {
		return "", InvalidUploadID{UploadID: uploadID}
	}
	// Calculate s3 compatible md5sum for complete multipart.
	s3MD5, err := completeMultipartMD5(parts...)
	if err != nil {
		return "", err
	}

	uploadIDPath := pathJoin(mpartMetaPrefix, bucket, object, uploadID)

	// Read metadata associated with the object from all disks.
	partsMetadata, errs := xl.readAllXLMetadata(minioMetaBucket, uploadIDPath)
	if err = xl.reduceError(errs); err != nil {
		return "", toObjectErr(err, minioMetaBucket, uploadIDPath)
	}

	// Calculate full object size.
	var objectSize int64

	// Pick one from the first valid metadata.
	xlMeta := pickValidXLMeta(partsMetadata)

	// Save current xl meta for validation.
	var currentXLMeta = xlMeta

	// Allocate parts similar to incoming slice.
	xlMeta.Parts = make([]objectPartInfo, len(parts))

	// Validate each part and then commit to disk.
	for i, part := range parts {
		partIdx := currentXLMeta.ObjectPartIndex(part.PartNumber)
		if partIdx == -1 {
			return "", InvalidPart{}
		}
		if currentXLMeta.Parts[partIdx].ETag != part.ETag {
			return "", BadDigest{}
		}
		// All parts except the last part has to be atleast 5MB.
		if (i < len(parts)-1) && !isMinAllowedPartSize(currentXLMeta.Parts[partIdx].Size) {
			return "", PartTooSmall{}
		}

		// Save for total object size.
		objectSize += currentXLMeta.Parts[partIdx].Size

		// Add incoming parts.
		xlMeta.Parts[i] = objectPartInfo{
			Number: part.PartNumber,
			ETag:   part.ETag,
			Size:   currentXLMeta.Parts[partIdx].Size,
			Name:   fmt.Sprintf("object%d", part.PartNumber),
		}
	}

	// Check if an object is present as one of the parent dir.
	if xl.parentDirIsObject(bucket, path.Dir(object)) {
		return "", toObjectErr(errFileAccessDenied, bucket, object)
	}

	// Save the final object size and modtime.
	xlMeta.Stat.Size = objectSize
	xlMeta.Stat.ModTime = time.Now().UTC()

	// Save successfully calculated md5sum.
	xlMeta.Meta["md5Sum"] = s3MD5
	uploadIDPath = path.Join(mpartMetaPrefix, bucket, object, uploadID)
	tempUploadIDPath := path.Join(tmpMetaPrefix, uploadID)

	// Update all xl metadata, make sure to not modify fields like
	// checksum which are different on each disks.
	for index := range partsMetadata {
		partsMetadata[index].Stat = xlMeta.Stat
		partsMetadata[index].Meta = xlMeta.Meta
		partsMetadata[index].Parts = xlMeta.Parts
	}

	// Write unique `xl.json` for each disk.
	if err = xl.writeUniqueXLMetadata(minioMetaBucket, tempUploadIDPath, partsMetadata); err != nil {
		return "", toObjectErr(err, minioMetaBucket, tempUploadIDPath)
	}
	rErr := xl.renameXLMetadata(minioMetaBucket, tempUploadIDPath, minioMetaBucket, uploadIDPath)
	if rErr != nil {
		return "", toObjectErr(rErr, minioMetaBucket, uploadIDPath)
	}
	// Hold write lock on the destination before rename
	nsMutex.Lock(bucket, object)
	defer nsMutex.Unlock(bucket, object)

	// Rename if an object already exists to temporary location.
	uniqueID := getUUID()
	err = xl.renameObject(bucket, object, minioMetaBucket, path.Join(tmpMetaPrefix, uniqueID))
	if err != nil {
		return "", toObjectErr(err, bucket, object)
	}

	// Remove parts that weren't present in CompleteMultipartUpload request
	for _, curpart := range currentXLMeta.Parts {
		if xlMeta.ObjectPartIndex(curpart.Number) == -1 {
			// Delete the missing part files. e.g,
			// Request 1: NewMultipart
			// Request 2: PutObjectPart 1
			// Request 3: PutObjectPart 2
			// Request 4: CompleteMultipartUpload --part 2
			// N.B. 1st part is not present. This part should be removed from the storage.
			xl.removeObjectPart(bucket, object, uploadID, curpart.Name)
		}
	}

	// Rename the multipart object to final location.
	if err = xl.renameObject(minioMetaBucket, uploadIDPath, bucket, object); err != nil {
		return "", toObjectErr(err, bucket, object)
	}

	// Delete the previously successfully renamed object.
	xl.deleteObject(minioMetaBucket, path.Join(tmpMetaPrefix, uniqueID))

	// Hold the lock so that two parallel complete-multipart-uploads do not
	// leave a stale uploads.json behind.
	nsMutex.Lock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object))
	defer nsMutex.Unlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object))

	// Validate if there are other incomplete upload-id's present for
	// the object, if yes do not attempt to delete 'uploads.json'.
	var disk StorageAPI
	var uploadsJSON uploadsV1
	for _, disk = range xl.getLoadBalancedQuorumDisks() {
		if disk == nil {
			continue
		}
		uploadsJSON, err = readUploadsJSON(bucket, object, disk)
		if err == errDiskNotFound {
			continue
		}
		break
	}
	if err != nil {
		return "", toObjectErr(err, minioMetaBucket, object)
	}
	// If we have successfully read `uploads.json`, then we proceed to
	// purge or update `uploads.json`.
	uploadIDIdx := uploadsJSON.Index(uploadID)
	if uploadIDIdx != -1 {
		uploadsJSON.Uploads = append(uploadsJSON.Uploads[:uploadIDIdx], uploadsJSON.Uploads[uploadIDIdx+1:]...)
	}
	if len(uploadsJSON.Uploads) > 0 {
		if err = updateUploadsJSON(bucket, object, uploadsJSON, xl.storageDisks...); err != nil {
			return "", toObjectErr(err, minioMetaBucket, path.Join(mpartMetaPrefix, bucket, object))
		}
		// Return success.
		return s3MD5, nil
	} // No more pending uploads for the object, proceed to delete
	// object completely from '.minio/multipart'.
	err = xl.deleteObject(minioMetaBucket, path.Join(mpartMetaPrefix, bucket, object))
	if err != nil {
		return "", toObjectErr(err, minioMetaBucket, path.Join(mpartMetaPrefix, bucket, object))
	}

	// Return md5sum.
	return s3MD5, nil
}

// abortMultipartUpload - wrapper for purging an ongoing multipart
// transaction, deletes uploadID entry from `uploads.json` and purges
// the directory at '.minio/multipart/bucket/object/uploadID' holding
// all the upload parts.
func (xl xlObjects) abortMultipartUpload(bucket, object, uploadID string) (err error) {
	// Cleanup all uploaded parts.
	if err = cleanupUploadedParts(bucket, object, uploadID, xl.storageDisks...); err != nil {
		return toObjectErr(err, bucket, object)
	}

	nsMutex.Lock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object))
	defer nsMutex.Unlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object))
	// Validate if there are other incomplete upload-id's present for
	// the object, if yes do not attempt to delete 'uploads.json'.
	var disk StorageAPI
	var uploadsJSON uploadsV1
	for _, disk = range xl.getLoadBalancedQuorumDisks() {
		if disk == nil {
			continue
		}
		uploadsJSON, err = readUploadsJSON(bucket, object, disk)
		if err == errDiskNotFound {
			continue
		}
		break
	}
	if err != nil {
		return toObjectErr(err, bucket, object)
	}
	uploadIDIdx := uploadsJSON.Index(uploadID)
	if uploadIDIdx != -1 {
		uploadsJSON.Uploads = append(uploadsJSON.Uploads[:uploadIDIdx], uploadsJSON.Uploads[uploadIDIdx+1:]...)
	}
	if len(uploadsJSON.Uploads) > 0 {
		// There are pending uploads for the same object, preserve
		// them update 'uploads.json' in-place.
		err = updateUploadsJSON(bucket, object, uploadsJSON, xl.storageDisks...)
		if err != nil {
			return toObjectErr(err, bucket, object)
		}
		return nil
	} // No more pending uploads for the object, we purge the entire
	// entry at '.minio/multipart/bucket/object'.
	if err = xl.deleteObject(minioMetaBucket, path.Join(mpartMetaPrefix, bucket, object)); err != nil {
		return toObjectErr(err, minioMetaBucket, path.Join(mpartMetaPrefix, bucket, object))
	}

	// Successfully purged.
	return nil
}

// AbortMultipartUpload - aborts an ongoing multipart operation
// signified by the input uploadID. This is an atomic operation
// doesn't require clients to initiate multiple such requests.
//
// All parts are purged from all disks and reference to the uploadID
// would be removed from the system, rollback is not possible on this
// operation.
//
// Implements S3 compatible Abort multipart API, slight difference is
// that this is an atomic idempotent operation. Subsequent calls have
// no affect and further requests to the same uploadID would not be honored.
func (xl xlObjects) AbortMultipartUpload(bucket, object, uploadID string) error {
	// Verify if bucket is valid.
	if !IsValidBucketName(bucket) {
		return BucketNameInvalid{Bucket: bucket}
	}
	if !xl.isBucketExist(bucket) {
		return BucketNotFound{Bucket: bucket}
	}
	if !IsValidObjectName(object) {
		return ObjectNameInvalid{Bucket: bucket, Object: object}
	}

	// Hold lock so that there is no competing complete-multipart-upload or put-object-part.
	nsMutex.Lock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object, uploadID))
	defer nsMutex.Unlock(minioMetaBucket, pathJoin(mpartMetaPrefix, bucket, object, uploadID))

	if !xl.isUploadIDExists(bucket, object, uploadID) {
		return InvalidUploadID{UploadID: uploadID}
	}
	err := xl.abortMultipartUpload(bucket, object, uploadID)
	return err
}
