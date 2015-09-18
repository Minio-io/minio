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

package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/minio/minio/pkg/auth"
	"github.com/minio/minio/pkg/donut"
	"github.com/minio/minio/pkg/probe"
)

const (
	authHeaderPrefix = "AWS4-HMAC-SHA256"
)

// StripAccessKeyID - strip only access key id from auth header
func StripAccessKeyID(ah string) (string, error) {
	if ah == "" {
		return "", errors.New("Missing auth header")
	}
	authFields := strings.Split(strings.TrimSpace(ah), ",")
	if len(authFields) != 3 {
		return "", errors.New("Missing fields in Auth header")
	}
	authPrefixFields := strings.Fields(authFields[0])
	if len(authPrefixFields) != 2 {
		return "", errors.New("Missing fields in Auth header")
	}
	if authPrefixFields[0] != authHeaderPrefix {
		return "", errors.New("Missing fields is Auth header")
	}
	credentials := strings.Split(strings.TrimSpace(authPrefixFields[1]), "=")
	if len(credentials) != 2 {
		return "", errors.New("Missing fields in Auth header")
	}
	if len(strings.Split(strings.TrimSpace(authFields[1]), "=")) != 2 {
		return "", errors.New("Missing fields in Auth header")
	}
	if len(strings.Split(strings.TrimSpace(authFields[2]), "=")) != 2 {
		return "", errors.New("Missing fields in Auth header")
	}
	accessKeyID := strings.Split(strings.TrimSpace(credentials[1]), "/")[0]
	if !auth.IsValidAccessKey(accessKeyID) {
		return "", errors.New("Invalid access key")
	}
	return accessKeyID, nil
}

// InitSignatureV4 initializing signature verification
func InitSignatureV4(req *http.Request) (*donut.Signature, *probe.Error) {
	authConfig, err := auth.LoadConfig()
	if err != nil {
		return nil, err.Trace()
	}

	if authConfig.NoAuth {
		return nil, nil
	}

	// strip auth from authorization header
	ah := req.Header.Get("Authorization")
	var accessKeyID string
	{
		var err error
		accessKeyID, err = StripAccessKeyID(ah)
		if err != nil {
			return nil, probe.NewError(err)
		}
	}
	if _, ok := authConfig.Users[accessKeyID]; !ok {
		return nil, probe.NewError(errors.New("AccessID not found"))
	}
	signature := &donut.Signature{
		AccessKeyID:     authConfig.Users[accessKeyID].AccessKeyID,
		SecretAccessKey: authConfig.Users[accessKeyID].SecretAccessKey,
		AuthHeader:      ah,
		Request:         req,
	}
	return signature, nil
}
