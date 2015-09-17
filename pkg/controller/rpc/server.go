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

package rpc

import (
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
)

// Server rpc server container
type Server struct {
	RPCServer *rpc.Server
}

// RegisterJSONCodec - register standard json codec
func (s Server) RegisterJSONCodec() {
	s.RPCServer.RegisterCodec(json.NewCodec(), "application/json")
}

// RegisterService - register new services
func (s Server) RegisterService(recv interface{}, name string) {
	s.RPCServer.RegisterService(recv, name)
}

// NewServer - provide a new instance of RPC server
func NewServer() *Server {
	s := &Server{}
	s.RPCServer = rpc.NewServer()
	return s
}

// ServeHTTP wrapper method for http.Handler interface
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.RPCServer.ServeHTTP(w, r)
}
