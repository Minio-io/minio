/*
 * Minimalist Object Storage, (C) 2015 Minio, Inc.
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

package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minio/minio/pkg/iodine"
	"github.com/minio/minio/pkg/storage/drivers"
	"github.com/minio/minio/pkg/utils/log"
)

func (server *minioAPI) isValidOp(w http.ResponseWriter, req *http.Request, acceptsContentType contentType) bool {
	vars := mux.Vars(req)
	bucket := vars["bucket"]

	bucketMetadata, err := server.driver.GetBucketMetadata(bucket)
	switch iodine.ToError(err).(type) {
	case drivers.BucketNotFound:
		{
			writeErrorResponse(w, req, NoSuchBucket, acceptsContentType, req.URL.Path)
			return false
		}
	case drivers.BucketNameInvalid:
		{
			writeErrorResponse(w, req, InvalidBucketName, acceptsContentType, req.URL.Path)
			return false
		}
	case nil:
		if _, err := stripAuth(req); err != nil {
			if bucketMetadata.ACL.IsPrivate() {
				writeErrorResponse(w, req, AccessDenied, acceptsContentType, req.URL.Path)
				return false
			}
			if bucketMetadata.ACL.IsPublicRead() && req.Method == "PUT" {
				writeErrorResponse(w, req, AccessDenied, acceptsContentType, req.URL.Path)
				return false
			}
		}
	default:
		{
			log.Error.Println(iodine.New(err, nil))
			writeErrorResponse(w, req, InternalError, acceptsContentType, req.URL.Path)
		}
	}
	return true
}

// GET Bucket (List Multipart uploads)
// -------------------------
// This operation lists in-progress multipart uploads. An in-progress
// multipart upload is a multipart upload that has been initiated,
// using the Initiate Multipart Upload request, but has not yet been completed or aborted.
// This operation returns at most 1,000 multipart uploads in the response.
//
func (server *minioAPI) listMultipartUploadsHandler(w http.ResponseWriter, req *http.Request) {
	acceptsContentType := getContentType(req)

	resources := getBucketMultipartResources(req.URL.Query())
	if resources.MaxUploads == 0 {
		resources.MaxUploads = maxObjectList
	}

	vars := mux.Vars(req)
	bucket := vars["bucket"]

	resources, err := server.driver.ListMultipartUploads(bucket, resources)
	switch iodine.ToError(err).(type) {
	case nil: // success
		{
			// generate response
			response := generateListMultipartUploadsResult(bucket, resources)
			encodedSuccessResponse := encodeSuccessResponse(response, acceptsContentType)
			// write headers
			setCommonHeaders(w, getContentTypeString(acceptsContentType), len(encodedSuccessResponse))
			// write body
			w.Write(encodedSuccessResponse)
		}
	case drivers.BucketNotFound:
		{
			writeErrorResponse(w, req, NoSuchBucket, acceptsContentType, req.URL.Path)
		}
	default:
		{
			log.Error.Println(iodine.New(err, nil))
			writeErrorResponse(w, req, InternalError, acceptsContentType, req.URL.Path)
		}
	}

}

// GET Bucket (List Objects)
// -------------------------
// This implementation of the GET operation returns some or all (up to 1000)
// of the objects in a bucket. You can use the request parameters as selection
// criteria to return a subset of the objects in a bucket.
//
func (server *minioAPI) listObjectsHandler(w http.ResponseWriter, req *http.Request) {
	acceptsContentType := getContentType(req)
	// verify if bucket allows this operation
	if !server.isValidOp(w, req, acceptsContentType) {
		return
	}

	if isRequestUploads(req.URL.Query()) {
		server.listMultipartUploadsHandler(w, req)
		return
	}

	resources := getBucketResources(req.URL.Query())
	if resources.Maxkeys == 0 {
		resources.Maxkeys = maxObjectList
	}

	vars := mux.Vars(req)
	bucket := vars["bucket"]

	objects, resources, err := server.driver.ListObjects(bucket, resources)
	switch iodine.ToError(err).(type) {
	case nil: // success
		{
			// generate response
			response := generateListObjectsResponse(bucket, objects, resources)
			encodedSuccessResponse := encodeSuccessResponse(response, acceptsContentType)
			// write headers
			setCommonHeaders(w, getContentTypeString(acceptsContentType), len(encodedSuccessResponse))
			// write body
			w.Write(encodedSuccessResponse)
		}
	case drivers.ObjectNotFound:
		{
			writeErrorResponse(w, req, NoSuchKey, acceptsContentType, req.URL.Path)
		}
	case drivers.ObjectNameInvalid:
		{
			writeErrorResponse(w, req, NoSuchKey, acceptsContentType, req.URL.Path)
		}
	default:
		{
			log.Error.Println(iodine.New(err, nil))
			writeErrorResponse(w, req, InternalError, acceptsContentType, req.URL.Path)
		}
	}
}

// GET Service
// -----------
// This implementation of the GET operation returns a list of all buckets
// owned by the authenticated sender of the request.
func (server *minioAPI) listBucketsHandler(w http.ResponseWriter, req *http.Request) {
	acceptsContentType := getContentType(req)
	// without access key credentials one cannot list buckets
	if _, err := stripAuth(req); err != nil {
		writeErrorResponse(w, req, AccessDenied, acceptsContentType, req.URL.Path)
		return
	}
	buckets, err := server.driver.ListBuckets()
	switch iodine.ToError(err).(type) {
	case nil:
		{
			// generate response
			response := generateListBucketsResponse(buckets)
			encodedSuccessResponse := encodeSuccessResponse(response, acceptsContentType)
			// write headers
			setCommonHeaders(w, getContentTypeString(acceptsContentType), len(encodedSuccessResponse))
			// write response
			w.Write(encodedSuccessResponse)
		}
	default:
		{
			log.Error.Println(iodine.New(err, nil))
			writeErrorResponse(w, req, InternalError, acceptsContentType, req.URL.Path)
		}
	}
}

// PUT Bucket
// ----------
// This implementation of the PUT operation creates a new bucket for authenticated request
func (server *minioAPI) putBucketHandler(w http.ResponseWriter, req *http.Request) {
	acceptsContentType := getContentType(req)
	// without access key credentials one cannot create a bucket
	if _, err := stripAuth(req); err != nil {
		writeErrorResponse(w, req, AccessDenied, acceptsContentType, req.URL.Path)
		return
	}
	if isRequestBucketACL(req.URL.Query()) {
		server.putBucketACLHandler(w, req)
		return
	}
	// read from 'x-amz-acl'
	aclType := getACLType(req)
	if aclType == unsupportedACLType {
		writeErrorResponse(w, req, NotImplemented, acceptsContentType, req.URL.Path)
		return
	}

	vars := mux.Vars(req)
	bucket := vars["bucket"]
	err := server.driver.CreateBucket(bucket, getACLTypeString(aclType))
	switch iodine.ToError(err).(type) {
	case nil:
		{
			// Make sure to add Location information here only for bucket
			w.Header().Set("Location", "/"+bucket)
			writeSuccessResponse(w, acceptsContentType)
		}
	case drivers.TooManyBuckets:
		{
			writeErrorResponse(w, req, TooManyBuckets, acceptsContentType, req.URL.Path)
		}
	case drivers.BucketNameInvalid:
		{
			writeErrorResponse(w, req, InvalidBucketName, acceptsContentType, req.URL.Path)
		}
	case drivers.BucketExists:
		{
			writeErrorResponse(w, req, BucketAlreadyExists, acceptsContentType, req.URL.Path)
		}
	default:
		{
			log.Error.Println(iodine.New(err, nil))
			writeErrorResponse(w, req, InternalError, acceptsContentType, req.URL.Path)
		}
	}
}

// PUT Bucket ACL
// ----------
// This implementation of the PUT operation modifies the bucketACL for authenticated request
func (server *minioAPI) putBucketACLHandler(w http.ResponseWriter, req *http.Request) {
	acceptsContentType := getContentType(req)
	// read from 'x-amz-acl'
	aclType := getACLType(req)
	if aclType == unsupportedACLType {
		writeErrorResponse(w, req, NotImplemented, acceptsContentType, req.URL.Path)
		return
	}

	vars := mux.Vars(req)
	bucket := vars["bucket"]
	err := server.driver.SetBucketMetadata(bucket, getACLTypeString(aclType))
	switch iodine.ToError(err).(type) {
	case nil:
		{
			writeSuccessResponse(w, acceptsContentType)
		}
	case drivers.BucketNameInvalid:
		{
			writeErrorResponse(w, req, InvalidBucketName, acceptsContentType, req.URL.Path)
		}
	case drivers.BucketNotFound:
		{
			writeErrorResponse(w, req, NoSuchBucket, acceptsContentType, req.URL.Path)
		}
	default:
		{
			log.Error.Println(iodine.New(err, nil))
			writeErrorResponse(w, req, InternalError, acceptsContentType, req.URL.Path)
		}
	}
}

// HEAD Bucket
// ----------
// This operation is useful to determine if a bucket exists.
// The operation returns a 200 OK if the bucket exists and you
// have permission to access it. Otherwise, the operation might
// return responses such as 404 Not Found and 403 Forbidden.
func (server *minioAPI) headBucketHandler(w http.ResponseWriter, req *http.Request) {
	acceptsContentType := getContentType(req)

	vars := mux.Vars(req)
	bucket := vars["bucket"]

	_, err := server.driver.GetBucketMetadata(bucket)
	switch iodine.ToError(err).(type) {
	case nil:
		{
			writeSuccessResponse(w, acceptsContentType)
		}
	case drivers.BucketNotFound:
		{
			error := getErrorCode(NoSuchBucket)
			w.WriteHeader(error.HTTPStatusCode)
		}
	case drivers.BucketNameInvalid:
		{
			error := getErrorCode(InvalidBucketName)
			w.WriteHeader(error.HTTPStatusCode)
		}
	default:
		{
			log.Error.Println(iodine.New(err, nil))
			error := getErrorCode(InternalError)
			w.WriteHeader(error.HTTPStatusCode)
		}
	}
}
