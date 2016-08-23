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
	"net/rpc"
	"testing"
)

// Wrapper for calling heal disk metadata rpc Handler
func TestControllerHandlerHealDiskMetadata(t *testing.T) {
	ExecObjectLayerTest(t, testHealDiskMetadataControllerHandler)
}

// testHealDiskMetadataControllerHandler - Test Heal Disk Metadata handler
func testHealDiskMetadataControllerHandler(obj ObjectLayer, instanceType string, t TestErrHandler) {
	// Register the API end points with XL/FS object layer.
	serverAddress, random, err := initTestControllerRPCEndPoint(obj)
	if err != nil {
		t.Fatal(err)
	}
	// initialize the server and obtain the credentials and root.
	// credentials are necessary to sign the HTTP request.
	rootPath, err := newTestConfig("us-east-1")
	if err != nil {
		t.Fatalf("Init Test config failed")
	}
	// remove the root folder after the test ends.
	defer removeAll(rootPath)

	client, err := rpc.DialHTTPPath("tcp", serverAddress, "/control"+random)
	if err != nil {
		t.Fatal("dialing:", err)
	}

	args := &HealDiskMetadataArgs{}
	reply := &HealDiskMetadataReply{}
	err = client.Call("Control.HealDiskMetadata", args, reply)
	if err != nil {
		t.Fatal("RPC Control.HealDiskMetadata call failed ", err)
	}
	if instanceType == "FS" && reply.Success {
		t.Errorf("Test should fail with FS")
	}
	if instanceType == "XL" && !reply.Success {
		t.Errorf("Test should succeed with XL")
	}
}
