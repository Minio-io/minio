/*
 * MinIO Cloud Storage, (C) 2018 MinIO, Inc.
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
	"bufio"
	"context"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"

	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/minio/minio/cmd/config"
	xhttp "github.com/minio/minio/cmd/http"
	xjwt "github.com/minio/minio/cmd/jwt"
	"github.com/minio/minio/cmd/logger"
)

var errDiskStale = errors.New("disk stale")

// To abstract a disk over network.
type storageRESTServer struct {
	storage *posix
}

func (s *storageRESTServer) writeErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(err.Error()))
	w.(http.Flusher).Flush()
}

// DefaultSkewTime - skew time is 15 minutes between minio peers.
const DefaultSkewTime = 15 * time.Minute

// Authenticates storage client's requests and validates for skewed time.
func storageServerRequestValidate(r *http.Request) error {
	token, err := jwtreq.AuthorizationHeaderExtractor.ExtractToken(r)
	if err != nil {
		if err == jwtreq.ErrNoTokenInRequest {
			return errNoAuthToken
		}
		return err
	}

	claims := xjwt.NewStandardClaims()
	if err = xjwt.ParseWithStandardClaims(token, claims, []byte(globalActiveCred.SecretKey)); err != nil {
		return errAuthentication
	}

	owner := claims.AccessKey == globalActiveCred.AccessKey || claims.Subject == globalActiveCred.AccessKey
	if !owner {
		return errAuthentication
	}

	if claims.Audience != r.URL.Query().Encode() {
		return errAuthentication
	}

	requestTimeStr := r.Header.Get("X-Minio-Time")
	requestTime, err := time.Parse(time.RFC3339, requestTimeStr)
	if err != nil {
		return err
	}
	utcNow := UTCNow()
	delta := requestTime.Sub(utcNow)
	if delta < 0 {
		delta *= -1
	}
	if delta > DefaultSkewTime {
		return fmt.Errorf("client time %v is too apart with server time %v", requestTime, utcNow)
	}

	return nil
}

// IsValid - To authenticate and verify the time difference.
func (s *storageRESTServer) IsValid(w http.ResponseWriter, r *http.Request) bool {
	if err := storageServerRequestValidate(r); err != nil {
		s.writeErrorResponse(w, err)
		return false
	}
	diskID := r.URL.Query().Get(storageRESTDiskID)
	if diskID == "" {
		// Request sent empty disk-id, we allow the request
		// as the peer might be coming up and trying to read format.json
		// or create format.json
		return true
	}
	storedDiskID, err := s.storage.GetDiskID()
	if err == nil && diskID == storedDiskID {
		// If format.json is available and request sent the right disk-id, we allow the request
		return true
	}
	s.writeErrorResponse(w, errDiskStale)
	return false
}

// DiskInfoHandler - returns disk info.
func (s *storageRESTServer) DiskInfoHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	info, err := s.storage.DiskInfo()
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	defer w.(http.Flusher).Flush()
	gob.NewEncoder(w).Encode(info)
}

func (s *storageRESTServer) UpdateBloomFilter(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	w.Header().Set(xhttp.ContentType, "text/event-stream")
	var req bloomFilterRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		logger.LogIf(r.Context(), err)
		s.writeErrorResponse(w, err)
		return
	}
	resp, err := cycleServerBloomFilter(r.Context(), req.Oldest, req.Current)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		logger.LogIf(r.Context(), err)
		s.writeErrorResponse(w, err)
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(resp)

	logger.LogIf(r.Context(), err)
}

func (s *storageRESTServer) CrawlAndGetDataUsageHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}

	w.Header().Set(xhttp.ContentType, "text/event-stream")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	var cache dataUsageCache
	err = cache.deserialize(b)
	if err != nil {
		logger.LogIf(r.Context(), err)
		s.writeErrorResponse(w, err)
		return
	}

	done := keepHTTPResponseAlive(w)
	usageInfo, err := s.storage.CrawlAndGetDataUsage(r.Context(), cache)
	done()

	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	w.Write(usageInfo.serialize())
	w.(http.Flusher).Flush()
}

// MakeVolHandler - make a volume.
func (s *storageRESTServer) MakeVolHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	err := s.storage.MakeVol(volume)
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// MakeVolBulkHandler - create multiple volumes as a bulk operation.
func (s *storageRESTServer) MakeVolBulkHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volumes := strings.Split(vars[storageRESTVolumes], ",")
	err := s.storage.MakeVolBulk(volumes...)
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// ListVolsHandler - list volumes.
func (s *storageRESTServer) ListVolsHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	infos, err := s.storage.ListVols()
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	gob.NewEncoder(w).Encode(&infos)
	w.(http.Flusher).Flush()
}

// StatVolHandler - stat a volume.
func (s *storageRESTServer) StatVolHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	info, err := s.storage.StatVol(volume)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	gob.NewEncoder(w).Encode(info)
	w.(http.Flusher).Flush()
}

// DeleteVolumeHandler - delete a volume.
func (s *storageRESTServer) DeleteVolHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	forceDelete := vars[storageRESTForceDelete] == "true"
	err := s.storage.DeleteVol(volume, forceDelete)
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// AppendFileHandler - append data from the request to the file specified.
func (s *storageRESTServer) AppendFileHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]

	buf := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, buf)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	err = s.storage.AppendFile(volume, filePath, buf)
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// CreateFileHandler - fallocate() space for a file and copy the contents from the request.
func (s *storageRESTServer) CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]

	fileSizeStr := vars[storageRESTLength]
	fileSize, err := strconv.Atoi(fileSizeStr)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	err = s.storage.CreateFile(volume, filePath, int64(fileSize), r.Body)
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// WriteAllHandler - write to file all content.
func (s *storageRESTServer) WriteAllHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]

	if r.ContentLength < 0 {
		s.writeErrorResponse(w, errInvalidArgument)
		return
	}

	err := s.storage.WriteAll(volume, filePath, io.LimitReader(r.Body, r.ContentLength))
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// StatFileHandler - stat a file.
func (s *storageRESTServer) StatFileHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]

	info, err := s.storage.StatFile(volume, filePath)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	gob.NewEncoder(w).Encode(info)
	w.(http.Flusher).Flush()
}

// ReadAllHandler - read all the contents of a file.
func (s *storageRESTServer) ReadAllHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]

	buf, err := s.storage.ReadAll(volume, filePath)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	w.Header().Set(xhttp.ContentLength, strconv.Itoa(len(buf)))
	w.Write(buf)
	w.(http.Flusher).Flush()
}

// ReadFileHandler - read section of a file.
func (s *storageRESTServer) ReadFileHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]
	offset, err := strconv.Atoi(vars[storageRESTOffset])
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	length, err := strconv.Atoi(vars[storageRESTLength])
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	if offset < 0 || length < 0 {
		s.writeErrorResponse(w, errInvalidArgument)
		return
	}
	var verifier *BitrotVerifier
	if vars[storageRESTBitrotAlgo] != "" {
		hashStr := vars[storageRESTBitrotHash]
		var hash []byte
		hash, err = hex.DecodeString(hashStr)
		if err != nil {
			s.writeErrorResponse(w, err)
			return
		}
		verifier = NewBitrotVerifier(BitrotAlgorithmFromString(vars[storageRESTBitrotAlgo]), hash)
	}
	buf := make([]byte, length)
	_, err = s.storage.ReadFile(volume, filePath, int64(offset), buf, verifier)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	w.Header().Set(xhttp.ContentLength, strconv.Itoa(len(buf)))
	w.Write(buf)
	w.(http.Flusher).Flush()
}

// ReadFileHandler - read section of a file.
func (s *storageRESTServer) ReadFileStreamHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]
	offset, err := strconv.Atoi(vars[storageRESTOffset])
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	length, err := strconv.Atoi(vars[storageRESTLength])
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}

	rc, err := s.storage.ReadFileStream(volume, filePath, int64(offset), int64(length))
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	defer rc.Close()

	w.Header().Set(xhttp.ContentLength, strconv.Itoa(length))

	io.Copy(w, rc)
	w.(http.Flusher).Flush()

}

// readMetadata func provides the function types for reading leaf metadata.
type readMetadataFunc func(buf []byte, volume, entry string) FileInfo

func readMetadata(buf []byte, volume, entry string) FileInfo {
	m, err := xlMetaV1UnmarshalJSON(context.Background(), buf)
	if err != nil {
		return FileInfo{}
	}
	return FileInfo{
		Volume:   volume,
		Name:     entry,
		ModTime:  m.Stat.ModTime,
		Size:     m.Stat.Size,
		Metadata: m.Meta,
		Parts:    m.Parts,
		Quorum:   m.Erasure.DataBlocks,
	}
}

// WalkHandler - remote caller to start walking at a requested directory path.
func (s *storageRESTServer) WalkSplunkHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	dirPath := vars[storageRESTDirPath]
	markerPath := vars[storageRESTMarkerPath]

	w.Header().Set(xhttp.ContentType, "text/event-stream")
	encoder := gob.NewEncoder(w)

	fch, err := s.storage.WalkSplunk(volume, dirPath, markerPath, r.Context().Done())
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	for fi := range fch {
		encoder.Encode(&fi)
	}
	w.(http.Flusher).Flush()
}

// WalkHandler - remote caller to start walking at a requested directory path.
func (s *storageRESTServer) WalkHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	dirPath := vars[storageRESTDirPath]
	markerPath := vars[storageRESTMarkerPath]
	recursive, err := strconv.ParseBool(vars[storageRESTRecursive])
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	leafFile := vars[storageRESTLeafFile]

	w.Header().Set(xhttp.ContentType, "text/event-stream")
	encoder := gob.NewEncoder(w)

	fch, err := s.storage.Walk(volume, dirPath, markerPath, recursive, leafFile, readMetadata, r.Context().Done())
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	for fi := range fch {
		encoder.Encode(&fi)
	}
	w.(http.Flusher).Flush()
}

// ListDirHandler - list a directory.
func (s *storageRESTServer) ListDirHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	dirPath := vars[storageRESTDirPath]
	leafFile := vars[storageRESTLeafFile]
	count, err := strconv.Atoi(vars[storageRESTCount])
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}

	entries, err := s.storage.ListDir(volume, dirPath, count, leafFile)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	gob.NewEncoder(w).Encode(&entries)
	w.(http.Flusher).Flush()
}

// DeleteFileHandler - delete a file.
func (s *storageRESTServer) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]

	err := s.storage.DeleteFile(volume, filePath)
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// DeleteFileBulkErrsResp - collection of deleteFile errors
// for bulk deletes
type DeleteFileBulkErrsResp struct {
	Errs []error
}

// DeleteFileBulkHandler - delete a file.
func (s *storageRESTServer) DeleteFileBulkHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := r.URL.Query()
	volume := vars.Get(storageRESTVolume)

	bio := bufio.NewScanner(r.Body)
	var filePaths []string
	for bio.Scan() {
		filePaths = append(filePaths, bio.Text())
	}

	if err := bio.Err(); err != nil {
		s.writeErrorResponse(w, err)
		return
	}

	dErrsResp := &DeleteFileBulkErrsResp{Errs: make([]error, len(filePaths))}

	w.Header().Set(xhttp.ContentType, "text/event-stream")
	encoder := gob.NewEncoder(w)
	done := keepHTTPResponseAlive(w)
	errs, err := s.storage.DeleteFileBulk(volume, filePaths)
	done()

	for idx := range filePaths {
		if err != nil {
			dErrsResp.Errs[idx] = StorageErr(err.Error())
		} else {
			if errs[idx] != nil {
				dErrsResp.Errs[idx] = StorageErr(errs[idx].Error())
			}
		}
	}

	encoder.Encode(dErrsResp)
	w.(http.Flusher).Flush()
}

// DeletePrefixesErrsResp - collection of delete errors
// for bulk prefixes deletes
type DeletePrefixesErrsResp struct {
	Errs []error
}

// DeletePrefixesHandler - delete a set of a prefixes.
func (s *storageRESTServer) DeletePrefixesHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := r.URL.Query()
	volume := vars.Get(storageRESTVolume)

	bio := bufio.NewScanner(r.Body)
	var prefixes []string
	for bio.Scan() {
		prefixes = append(prefixes, bio.Text())
	}

	if err := bio.Err(); err != nil {
		s.writeErrorResponse(w, err)
		return
	}

	dErrsResp := &DeletePrefixesErrsResp{Errs: make([]error, len(prefixes))}

	w.Header().Set(xhttp.ContentType, "text/event-stream")
	encoder := gob.NewEncoder(w)
	done := keepHTTPResponseAlive(w)
	errs, err := s.storage.DeletePrefixes(volume, prefixes)
	done()
	for idx := range prefixes {
		if err != nil {
			dErrsResp.Errs[idx] = StorageErr(err.Error())
		} else {
			if errs[idx] != nil {
				dErrsResp.Errs[idx] = StorageErr(errs[idx].Error())
			}
		}
	}
	encoder.Encode(dErrsResp)
	w.(http.Flusher).Flush()
}

// RenameFileHandler - rename a file.
func (s *storageRESTServer) RenameFileHandler(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	srcVolume := vars[storageRESTSrcVolume]
	srcFilePath := vars[storageRESTSrcPath]
	dstVolume := vars[storageRESTDstVolume]
	dstFilePath := vars[storageRESTDstPath]
	err := s.storage.RenameFile(srcVolume, srcFilePath, dstVolume, dstFilePath)
	if err != nil {
		s.writeErrorResponse(w, err)
	}
}

// keepHTTPResponseAlive can be used to avoid timeouts with long storage
// operations, such as bitrot verification or data usage crawling.
// Every 10 seconds a space character is sent.
// The returned function should always be called to release resources.
// waitForHTTPResponse should be used to the receiving side.
func keepHTTPResponseAlive(w http.ResponseWriter) func() {
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				w.Write([]byte(" "))
				w.(http.Flusher).Flush()
			case doneCh <- struct{}{}:
				w.Write([]byte{0})
				ticker.Stop()
				return
			}
		}
	}()
	return func() {
		// Indicate we are ready to write.
		<-doneCh
		// Wait for channel to be closed so we don't race on writes.
		<-doneCh
	}
}

// waitForHTTPResponse will wait for responses where keepHTTPResponseAlive
// has been used.
// The returned reader contains the payload.
func waitForHTTPResponse(respBody io.Reader) (io.Reader, error) {
	reader := bufio.NewReader(respBody)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b != ' ' {
			if b != 0 {
				reader.UnreadByte()
			}
			break
		}
	}
	return reader, nil
}

// VerifyFileResp - VerifyFile()'s response.
type VerifyFileResp struct {
	Err error
}

// VerifyFile - Verify the file for bitrot errors.
func (s *storageRESTServer) VerifyFile(w http.ResponseWriter, r *http.Request) {
	if !s.IsValid(w, r) {
		return
	}
	vars := mux.Vars(r)
	volume := vars[storageRESTVolume]
	filePath := vars[storageRESTFilePath]
	size, err := strconv.ParseInt(vars[storageRESTLength], 10, 0)
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	shardSize, err := strconv.Atoi(vars[storageRESTShardSize])
	if err != nil {
		s.writeErrorResponse(w, err)
		return
	}
	hashStr := vars[storageRESTBitrotHash]
	var hash []byte
	if hashStr != "" {
		hash, err = hex.DecodeString(hashStr)
		if err != nil {
			s.writeErrorResponse(w, err)
			return
		}
	}
	algoStr := vars[storageRESTBitrotAlgo]
	if algoStr == "" {
		s.writeErrorResponse(w, errInvalidArgument)
		return
	}
	w.Header().Set(xhttp.ContentType, "text/event-stream")
	encoder := gob.NewEncoder(w)
	done := keepHTTPResponseAlive(w)
	err = s.storage.VerifyFile(volume, filePath, size, BitrotAlgorithmFromString(algoStr), hash, int64(shardSize))
	done()
	vresp := &VerifyFileResp{}
	if err != nil {
		vresp.Err = StorageErr(err.Error())
	}
	encoder.Encode(vresp)
	w.(http.Flusher).Flush()
}

// registerStorageRPCRouter - register storage rpc router.
func registerStorageRESTHandlers(router *mux.Router, endpointZones EndpointZones) {
	for _, ep := range endpointZones {
		for _, endpoint := range ep.Endpoints {
			if !endpoint.IsLocal {
				continue
			}
			storage, err := newPosix(endpoint.Path)
			if err != nil {
				// Show a descriptive error with a hint about how to fix it.
				var username string
				if u, err := user.Current(); err == nil {
					username = u.Username
				} else {
					username = "<your-username>"
				}
				hint := fmt.Sprintf("Run the following command to add the convenient permissions: `sudo chown %s %s && sudo chmod u+rxw %s`",
					username, endpoint.Path, endpoint.Path)
				logger.Fatal(config.ErrUnableToWriteInBackend(err).Hint(hint),
					"Unable to initialize posix backend")
			}

			server := &storageRESTServer{storage: storage}

			subrouter := router.PathPrefix(path.Join(storageRESTPrefix, endpoint.Path)).Subrouter()

			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodDiskInfo).HandlerFunc(httpTraceHdrs(server.DiskInfoHandler))
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodCrawlAndGetDataUsage).HandlerFunc(httpTraceHdrs(server.CrawlAndGetDataUsageHandler))
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodUpdateBloomFilter).HandlerFunc(httpTraceHdrs(server.UpdateBloomFilter))
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodMakeVol).HandlerFunc(httpTraceHdrs(server.MakeVolHandler)).Queries(restQueries(storageRESTVolume)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodMakeVolBulk).HandlerFunc(httpTraceHdrs(server.MakeVolBulkHandler)).Queries(restQueries(storageRESTVolumes)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodStatVol).HandlerFunc(httpTraceHdrs(server.StatVolHandler)).Queries(restQueries(storageRESTVolume)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodDeleteVol).HandlerFunc(httpTraceHdrs(server.DeleteVolHandler)).Queries(restQueries(storageRESTVolume)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodListVols).HandlerFunc(httpTraceHdrs(server.ListVolsHandler))

			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodAppendFile).HandlerFunc(httpTraceHdrs(server.AppendFileHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodWriteAll).HandlerFunc(httpTraceHdrs(server.WriteAllHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodCreateFile).HandlerFunc(httpTraceHdrs(server.CreateFileHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath, storageRESTLength)...)

			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodStatFile).HandlerFunc(httpTraceHdrs(server.StatFileHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodReadAll).HandlerFunc(httpTraceHdrs(server.ReadAllHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodReadFile).HandlerFunc(httpTraceHdrs(server.ReadFileHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath, storageRESTOffset, storageRESTLength, storageRESTBitrotAlgo, storageRESTBitrotHash)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodReadFileStream).HandlerFunc(httpTraceHdrs(server.ReadFileStreamHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath, storageRESTOffset, storageRESTLength)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodListDir).HandlerFunc(httpTraceHdrs(server.ListDirHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTDirPath, storageRESTCount, storageRESTLeafFile)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodWalk).HandlerFunc(httpTraceHdrs(server.WalkHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTDirPath, storageRESTMarkerPath, storageRESTRecursive, storageRESTLeafFile)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodWalkSplunk).HandlerFunc(httpTraceHdrs(server.WalkSplunkHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTDirPath, storageRESTMarkerPath)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodDeletePrefixes).HandlerFunc(httpTraceHdrs(server.DeletePrefixesHandler)).
				Queries(restQueries(storageRESTVolume)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodDeleteFile).HandlerFunc(httpTraceHdrs(server.DeleteFileHandler)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodDeleteFileBulk).HandlerFunc(httpTraceHdrs(server.DeleteFileBulkHandler)).
				Queries(restQueries(storageRESTVolume)...)

			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodRenameFile).HandlerFunc(httpTraceHdrs(server.RenameFileHandler)).
				Queries(restQueries(storageRESTSrcVolume, storageRESTSrcPath, storageRESTDstVolume, storageRESTDstPath)...)
			subrouter.Methods(http.MethodPost).Path(storageRESTVersionPrefix + storageRESTMethodVerifyFile).HandlerFunc(httpTraceHdrs(server.VerifyFile)).
				Queries(restQueries(storageRESTVolume, storageRESTFilePath, storageRESTBitrotAlgo, storageRESTBitrotHash, storageRESTLength, storageRESTShardSize)...)
		}
	}
}
