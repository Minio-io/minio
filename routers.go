/*
 * Minio Cloud Storage, (C) 2015, 2016 Minio, Inc.
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
	"net/http"

	router "github.com/gorilla/mux"
	"github.com/minio/minio/pkg/fs"
)

// configureServer handler returns final handler for the http server.
func configureServerHandler(filesystem fs.Filesystem) http.Handler {
	// Initialize API.
	api := storageAPI{
		Filesystem: filesystem,
	}
	// Initialize Web.
	web := &webAPI{
		Filesystem: filesystem,
	}

	// Initialize router.
	mux := router.NewRouter()

	// Register all routers.
	registerWebRouter(mux, web)
	registerAPIRouter(mux, api)
	// Add new routers here.

	// List of some generic handlers which are applied for all
	// incoming requests.
	var handlerFns = []HandlerFunc{
		// Redirect some pre-defined browser request paths to a static
		// location prefix.
		setBrowserRedirectHandler,
		// Validates if incoming request is for restricted buckets.
		setPrivateBucketHandler,
		// Adds cache control for all browser requests.
		setBrowserCacheControlHandler,
		// Validates all incoming requests to have a valid date header.
		setTimeValidityHandler,
		// CORS setting for all browser API requests.
		setCorsHandler,
		// Validates all incoming URL resources, for invalid/unsupported
		// resources client receives a HTTP error.
		setIgnoreResourcesHandler,
		// Auth handler verifies incoming authorization headers and
		// routes them accordingly. Client receives a HTTP error for
		// invalid/unsupported signatures.
		setAuthHandler,
		// Add new handlers here.
	}

	// Register rest of the handlers.
	return registerHandlers(mux, handlerFns...)
}
