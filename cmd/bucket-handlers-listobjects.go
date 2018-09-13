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

package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio/pkg/policy"
)

// Validate all the ListObjects query arguments, returns an APIErrorCode
// if one of the args do not meet the required conditions.
// Special conditions required by Minio server are as below
// - delimiter if set should be equal to '/', otherwise the request is rejected.
// - marker if set should have a common prefix with 'prefix' param, otherwise
//   the request is rejected.
func validateListObjectsArgs(prefix, marker, delimiter string, maxKeys int) APIErrorCode {
	// Max keys cannot be negative.
	if maxKeys < 0 {
		return ErrInvalidMaxKeys
	}

	/// Minio special conditions for ListObjects.

	// Verify if delimiter is anything other than '/', which we do not support.
	if delimiter != "" && delimiter != "/" {
		return ErrNotImplemented
	}
	// Success.
	return ErrNone
}

// Validate all the ListObjectVersions query arguments, returns an APIErrorCode
// if one of the args do not meet the required conditions.
// Special conditions required by Minio server are as below
// - delimiter if set should be equal to '/', otherwise the request is rejected.
// - marker if set should have a common prefix with 'prefix' param, otherwise
//   the request is rejected.
func validateListObjectsVersionsArgs(delimiter, encodingType, keyMarker string, maxKeys int, prefix, versionIDMarker string) APIErrorCode {
	// Max keys cannot be negative.
	if maxKeys < 0 {
		return ErrInvalidMaxKeys
	}

	/// Minio special conditions for ListObjects.

	// Verify if delimiter is anything other than '/', which we do not support.
	if delimiter != "" && delimiter != "/" {
		return ErrNotImplemented
	}
	// Success.
	return ErrNone
}

// ListObjectsVersionsHandler - GET Bucket Object versions
func (api objectAPIHandlers) ListObjectsVersionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r, w, "ListObjectVersions")

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	objectAPI := api.ObjectAPI()
	if objectAPI == nil {
		writeErrorResponse(w, ErrServerNotInitialized, r.URL)
		return
	}

	if s3Error := checkRequestAuthType(ctx, r, policy.ListBucketAction, bucket, ""); s3Error != ErrNone {
		writeErrorResponse(w, s3Error, r.URL)
		return
	}

	urlValues := r.URL.Query()

	// Extract all the listObjectsVersions query params to their native values.
	delimiter, encodingType, keyMarker, maxkeys, prefix, versionIDMarker := getListObjectsVersionsArgs(urlValues)

	// Validate the query params before beginning to serve the request.
	// fetch-owner is not validated since it is a boolean
	if s3Error := validateListObjectsVersionsArgs(delimiter, encodingType, keyMarker, maxkeys, prefix, versionIDMarker); s3Error != ErrNone {
		writeErrorResponse(w, s3Error, r.URL)
		return
	}

	listObjectsVersions := objectAPI.ListObjectVersions

	// FIXME: shall we cache listing objects versions
	// if api.CacheAPI() != nil {
	//	listObjectsV2 = api.CacheAPI().ListObjectsV2
	//}

	// Initiate a list objects operation based on the input params.
	// On success would return back ListObjectsInfo object to be
	// marshaled into S3 compatible XML header.
	listObjectsVersionsInfo, err := listObjectsVersions(ctx, bucket, prefix, delimiter, keyMarker, versionIDMarker, maxkeys)
	if err != nil {
		writeErrorResponse(w, toAPIErrorCode(err), r.URL)
		return
	}

	// FIXME: check if size in versions differentiate between encryption/plain state of object
	/* for i := range listObjectsVersionsInfo.Versions {
		if listObjectsVersionsInfo.Versions[i].IsEncrypted() {
			listObjectsVersionsInfo.Versions[i].Size, err = listObjectsVersionsInfo.Versions[i].DecryptedSize()
			if err != nil {
				writeErrorResponse(w, toAPIErrorCode(err), r.URL)
				return
			}
		}
	} */

	response := generateListObjectsVersionsResponse(bucket, prefix, delimiter, listObjectsVersionsInfo.IsTruncated, maxkeys, listObjectsVersionsInfo.Versions, listObjectsVersionsInfo.DeleteMarkers)

	// Write success response.
	writeSuccessResponseXML(w, encodeResponse(response))
}

// ListObjectsV2Handler - GET Bucket (List Objects) Version 2.
// --------------------------
// This implementation of the GET operation returns some or all (up to 1000)
// of the objects in a bucket. You can use the request parameters as selection
// criteria to return a subset of the objects in a bucket.
//
// NOTE: It is recommended that this API to be used for application development.
// Minio continues to support ListObjectsV1 for supporting legacy tools.
func (api objectAPIHandlers) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r, w, "ListObjectsV2")

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	objectAPI := api.ObjectAPI()
	if objectAPI == nil {
		writeErrorResponse(w, ErrServerNotInitialized, r.URL)
		return
	}

	if s3Error := checkRequestAuthType(ctx, r, policy.ListBucketAction, bucket, ""); s3Error != ErrNone {
		writeErrorResponse(w, s3Error, r.URL)
		return
	}

	urlValues := r.URL.Query()

	// Extract all the listObjectsV2 query params to their native values.
	prefix, token, startAfter, delimiter, fetchOwner, maxKeys, _, errCode := getListObjectsV2Args(urlValues)

	if errCode != ErrNone {
		writeErrorResponse(w, errCode, r.URL)
		return
	}

	// Validate the query params before beginning to serve the request.
	// fetch-owner is not validated since it is a boolean
	if s3Error := validateListObjectsArgs(prefix, token, delimiter, maxKeys); s3Error != ErrNone {
		writeErrorResponse(w, s3Error, r.URL)
		return
	}
	listObjectsV2 := objectAPI.ListObjectsV2
	if api.CacheAPI() != nil {
		listObjectsV2 = api.CacheAPI().ListObjectsV2
	}
	// Inititate a list objects operation based on the input params.
	// On success would return back ListObjectsInfo object to be
	// marshaled into S3 compatible XML header.
	listObjectsV2Info, err := listObjectsV2(ctx, bucket, prefix, token, delimiter, maxKeys, fetchOwner, startAfter)
	if err != nil {
		writeErrorResponse(w, toAPIErrorCode(err), r.URL)
		return
	}

	for i := range listObjectsV2Info.Objects {
		if listObjectsV2Info.Objects[i].IsEncrypted() {
			listObjectsV2Info.Objects[i].Size, err = listObjectsV2Info.Objects[i].DecryptedSize()
			if err != nil {
				writeErrorResponse(w, toAPIErrorCode(err), r.URL)
				return
			}
		}
	}

	response := generateListObjectsV2Response(bucket, prefix, token, listObjectsV2Info.NextContinuationToken, startAfter,
		delimiter, fetchOwner, listObjectsV2Info.IsTruncated, maxKeys, listObjectsV2Info.Objects, listObjectsV2Info.Prefixes)

	// Write success response.
	writeSuccessResponseXML(w, encodeResponse(response))
}

// ListObjectsV1Handler - GET Bucket (List Objects) Version 1.
// --------------------------
// This implementation of the GET operation returns some or all (up to 1000)
// of the objects in a bucket. You can use the request parameters as selection
// criteria to return a subset of the objects in a bucket.
//
func (api objectAPIHandlers) ListObjectsV1Handler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r, w, "ListObjectsV1")

	vars := mux.Vars(r)
	bucket := vars["bucket"]

	objectAPI := api.ObjectAPI()
	if objectAPI == nil {
		writeErrorResponse(w, ErrServerNotInitialized, r.URL)
		return
	}

	if s3Error := checkRequestAuthType(ctx, r, policy.ListBucketAction, bucket, ""); s3Error != ErrNone {
		writeErrorResponse(w, s3Error, r.URL)
		return
	}

	// Extract all the listObjectsV1 query params to their native values.
	prefix, marker, delimiter, maxKeys, _ := getListObjectsV1Args(r.URL.Query())

	// Validate the maxKeys lowerbound. When maxKeys > 1000, S3 returns 1000 but
	// does not throw an error.
	if maxKeys < 0 {
		writeErrorResponse(w, ErrInvalidMaxKeys, r.URL)
		return
	} // Validate all the query params before beginning to serve the request.
	if s3Error := validateListObjectsArgs(prefix, marker, delimiter, maxKeys); s3Error != ErrNone {
		writeErrorResponse(w, s3Error, r.URL)
		return
	}
	listObjects := objectAPI.ListObjects
	if api.CacheAPI() != nil {
		listObjects = api.CacheAPI().ListObjects
	}
	// Inititate a list objects operation based on the input params.
	// On success would return back ListObjectsInfo object to be
	// marshaled into S3 compatible XML header.
	listObjectsInfo, err := listObjects(ctx, bucket, prefix, marker, delimiter, maxKeys)
	if err != nil {
		writeErrorResponse(w, toAPIErrorCode(err), r.URL)
		return
	}

	for i := range listObjectsInfo.Objects {
		if listObjectsInfo.Objects[i].IsEncrypted() {
			listObjectsInfo.Objects[i].Size, err = listObjectsInfo.Objects[i].DecryptedSize()
			if err != nil {
				writeErrorResponse(w, toAPIErrorCode(err), r.URL)
				return
			}
		}
	}

	response := generateListObjectsV1Response(bucket, prefix, marker, delimiter, maxKeys, listObjectsInfo)

	// Write success response.
	writeSuccessResponseXML(w, encodeResponse(response))
}
