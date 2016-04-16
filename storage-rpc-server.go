package main

import (
	"io"
	"net/http"
	"net/rpc"
	"os"
	"strconv"

	router "github.com/gorilla/mux"
)

// Storage server implements rpc primitives to facilitate exporting a
// disk over a network.
type storageServer struct {
	storage StorageAPI
}

/// Volume operations handlers

// MakeVolHandler - make vol handler is rpc wrapper for MakeVol operation.
func (s *storageServer) MakeVolHandler(arg *string, reply *GenericReply) error {
	return s.storage.MakeVol(*arg)
}

// ListVolsHandler - list vols handler is rpc wrapper for ListVols operation.
func (s *storageServer) ListVolsHandler(arg *string, reply *ListVolsReply) error {
	vols, err := s.storage.ListVols()
	if err != nil {
		return err
	}
	reply.Vols = vols
	return nil
}

// StatVolHandler - stat vol handler is a rpc wrapper for StatVol operation.
func (s *storageServer) StatVolHandler(arg *string, reply *VolInfo) error {
	volInfo, err := s.storage.StatVol(*arg)
	if err != nil {
		return err
	}
	*reply = volInfo
	return nil
}

// DeleteVolHandler - delete vol handler is a rpc wrapper for
// DeleteVol operation.
func (s *storageServer) DeleteVolHandler(arg *string, reply *GenericReply) error {
	return s.storage.DeleteVol(*arg)
}

/// File operations

// ListFilesHandler - list files handler.
func (s *storageServer) ListFilesHandler(arg *ListFilesArgs, reply *ListFilesReply) error {
	files, eof, err := s.storage.ListFiles(arg.Vol, arg.Prefix, arg.Marker, arg.Recursive, arg.Count)
	if err != nil {
		return err
	}
	reply.Files = files
	reply.EOF = eof
	return nil
}

// StatFileHandler - stat file handler is rpc wrapper to stat file.
func (s *storageServer) StatFileHandler(arg *StatFileArgs, reply *FileInfo) error {
	fileInfo, err := s.storage.StatFile(arg.Vol, arg.Path)
	if err != nil {
		return err
	}
	*reply = fileInfo
	return nil
}

// DeleteFileHandler - delete file handler is rpc wrapper to delete file.
func (s *storageServer) DeleteFileHandler(arg *DeleteFileArgs, reply *GenericReply) error {
	return s.storage.DeleteFile(arg.Vol, arg.Path)
}

// registerStorageRPCRouter - register storage rpc router.
func registerStorageRPCRouter(mux *router.Router, storageAPI StorageAPI) {
	stServer := &storageServer{
		storage: storageAPI,
	}
	storageRPCServer := rpc.NewServer()
	storageRPCServer.RegisterName("Storage", stServer)
	storageRouter := mux.NewRoute().PathPrefix(reservedBucket).Subrouter()
	// Add minio storage routes.
	storageRouter.Path("/rpc/storage").Handler(storageRPCServer)
	// StreamUpload - stream upload handler.
	storageRouter.Methods("POST").Path("/rpc/storage/upload/{volume}/{path:.+}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := router.Vars(r)
		volume := vars["volume"]
		path := vars["path"]
		writeCloser, err := stServer.storage.CreateFile(volume, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader := r.Body
		if _, err = io.Copy(writeCloser, reader); err != nil {
			safeCloseAndRemove(writeCloser)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeCloser.Close()
		reader.Close()
	})
	// StreamDownloadHandler - stream download handler.
	storageRouter.Methods("GET").Path("/rpc/storage/download/{volume}/{path:.+}").Queries("offset", "{offset:.*}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := router.Vars(r)
		volume := vars["volume"]
		path := vars["path"]
		offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		readCloser, err := stServer.storage.ReadFile(volume, path, offset)
		if err != nil {
			httpErr := http.StatusBadRequest
			if os.IsNotExist(err) {
				httpErr = http.StatusNotFound
			}
			http.Error(w, err.Error(), httpErr)
			return
		}
		io.Copy(w, readCloser)
		w.(http.Flusher).Flush()
		readCloser.Close()
	})
}
