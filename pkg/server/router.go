/*
 * Minio Cloud Storage, (C) 2014 Minio, Inc.
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

package server

import (
	"net/http"

	router "github.com/gorilla/mux"
	"github.com/minio/minio/pkg/rpc"
	"github.com/minio/minio/pkg/server/api"
)

// registerAPI - register all the object API handlers to their respective paths
func registerAPI(mux *router.Router, a api.Minio) {
	mux.HandleFunc("/", a.ListBucketsHandler).Methods("GET")
	mux.HandleFunc("/{bucket}", a.ListObjectsHandler).Methods("GET")
	mux.HandleFunc("/{bucket}", a.PutBucketHandler).Methods("PUT")
	mux.HandleFunc("/{bucket}", a.HeadBucketHandler).Methods("HEAD")
	mux.HandleFunc("/{bucket}/{object:.*}", a.HeadObjectHandler).Methods("HEAD")
	mux.HandleFunc("/{bucket}/{object:.*}", a.PutObjectPartHandler).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}").Methods("PUT")
	mux.HandleFunc("/{bucket}/{object:.*}", a.ListObjectPartsHandler).Queries("uploadId", "{uploadId:.*}").Methods("GET")
	mux.HandleFunc("/{bucket}/{object:.*}", a.CompleteMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}").Methods("POST")
	mux.HandleFunc("/{bucket}/{object:.*}", a.NewMultipartUploadHandler).Methods("POST")
	mux.HandleFunc("/{bucket}/{object:.*}", a.AbortMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}").Methods("DELETE")
	mux.HandleFunc("/{bucket}/{object:.*}", a.GetObjectHandler).Methods("GET")
	mux.HandleFunc("/{bucket}/{object:.*}", a.PutObjectHandler).Methods("PUT")

	// not implemented yet
	mux.HandleFunc("/{bucket}", a.DeleteBucketHandler).Methods("DELETE")

	// unsupported API
	mux.HandleFunc("/{bucket}/{object:.*}", a.DeleteObjectHandler).Methods("DELETE")
}

func registerCustomMiddleware(mux *router.Router, mwHandlers ...api.MiddlewareHandler) http.Handler {
	var f http.Handler
	f = mux
	for _, mw := range mwHandlers {
		f = mw(f)
	}
	return f
}

// getAPIHandler api handler
func getAPIHandler(conf api.Config) (http.Handler, api.Minio) {
	var mwHandlers = []api.MiddlewareHandler{
		api.ValidContentTypeHandler,
		api.TimeValidityHandler,
		api.IgnoreResourcesHandler,
		api.ValidateAuthHeaderHandler,
		// api.LoggingHandler, // Disabled logging until we bring in external logging support
		api.CorsHandler,
	}

	mux := router.NewRouter()
	minioAPI := api.New()
	registerAPI(mux, minioAPI)
	apiHandler := registerCustomMiddleware(mux, mwHandlers...)
	return apiHandler, minioAPI
}

// getRPCHandler rpc handler
func getRPCHandler() http.Handler {
	s := rpc.NewServer()
	s.RegisterJSONCodec()
	// Add new RPC services here
	return registerRPC(router.NewRouter(), s)
}

// registerRPC - register rpc handlers
func registerRPC(mux *router.Router, s *rpc.Server) http.Handler {
	mux.Handle("/rpc", s)
	return mux
}
