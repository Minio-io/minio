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

import "testing"

func testAuthenticate(authType string, t *testing.T) {
	testPath, err := newTestConfig("us-east-1")
	if err != nil {
		t.Fatalf("unable initialize config file, %s", err)
	}
	defer removeAll(testPath)

	serverCred := serverConfig.GetCredential()

	// Define test cases.
	testCases := []struct {
		accessKey   string
		secretKey   string
		expectedErr error
	}{
		// Access key too small.
		{"user", "pass", errInvalidAccessKeyLength},
		// Access key too long.
		{"user12345678901234567", "pass", errInvalidAccessKeyLength},
		// Access key contains unsuppported characters.
		{"!@#$", "pass", errInvalidAccessKeyLength},
		// Secret key too small.
		{"myuser", "pass", errInvalidSecretKeyLength},
		// Secret key too long.
		{"myuser", "pass1234567890123456789012345678901234567", errInvalidSecretKeyLength},
		// Authentication error.
		{"myuser", "mypassword", errInvalidAccessKeyID},
		// Authentication error.
		{serverCred.AccessKey, "mypassword", errAuthentication},
		// Success.
		{serverCred.AccessKey, serverCred.SecretKey, nil},
		// Success when access key contains leading/trailing spaces.
		{" " + serverCred.AccessKey + " ", serverCred.SecretKey, nil},
	}

	// Run tests.
	for _, testCase := range testCases {
		var err error
		if authType == "node" {
			_, err = authenticateNode(testCase.accessKey, testCase.secretKey)
		} else if authType == "web" {
			_, err = authenticateWeb(testCase.accessKey, testCase.secretKey)
		}

		if testCase.expectedErr != nil {
			if err == nil {
				t.Fatalf("%+v: expected: %s, got: <nil>", testCase, testCase.expectedErr)
			}
			if testCase.expectedErr.Error() != err.Error() {
				t.Fatalf("%+v: expected: %s, got: %s", testCase, testCase.expectedErr, err)
			}
		} else if err != nil {
			t.Fatalf("%+v: expected: <nil>, got: %s", testCase, err)
		}
	}
}

func TestNodeAuthenticate(t *testing.T) {
	testAuthenticate("node", t)
}

func TestWebAuthenticate(t *testing.T) {
	testAuthenticate("web", t)
}
