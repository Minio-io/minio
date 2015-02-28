/*
 * Mini Object Storage, (C) 2015 Minio, Inc.
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

package v1

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"log"
	"testing"

	"github.com/minio-io/minio/pkg/utils/checksum/crc32c"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSingleWrite(c *C) {
	//var b io.ReadWriteSeeker
	var testBuffer bytes.Buffer

	testData := "Hello, World"
	testLength := uint64(len(testData))
	err := Write(&testBuffer, bytes.NewBufferString(testData), testLength)
	c.Assert(err, IsNil)

	testBufferLength := uint64(testBuffer.Len())
	log.Println(testBuffer.Bytes())

	// we test our crc here too
	headerBytes := testBuffer.Bytes()[0:28]
	expectedCrc, err := crc32c.Crc32c(headerBytes)
	c.Assert(err, IsNil)

	// magic mini
	magicMini := make([]byte, 4)
	testBuffer.Read(magicMini)
	c.Assert(magicMini, DeepEquals, []byte{'M', 'I', 'N', 'I'})

	// major version
	majorVersion := make([]byte, 2)
	testBuffer.Read(majorVersion)
	c.Assert(binary.LittleEndian.Uint16(majorVersion), DeepEquals, uint16(1))

	// minor version
	minorVersion := make([]byte, 2)
	testBuffer.Read(minorVersion)
	c.Assert(binary.LittleEndian.Uint16(minorVersion), DeepEquals, uint16(0))

	// patch version
	patchVersion := make([]byte, 2)
	testBuffer.Read(patchVersion)
	c.Assert(binary.LittleEndian.Uint16(patchVersion), DeepEquals, uint16(0))

	// reserved version
	reservedVersion := make([]byte, 2)
	testBuffer.Read(reservedVersion)
	c.Assert(binary.LittleEndian.Uint16(reservedVersion), DeepEquals, uint16(0))

	// reserved
	reserved := make([]byte, 8)
	testBuffer.Read(reserved)
	c.Assert(binary.LittleEndian.Uint64(reserved), DeepEquals, uint64(0))

	// data length
	length := make([]byte, 8)
	testBuffer.Read(length)
	c.Assert(binary.LittleEndian.Uint64(length), DeepEquals, testLength)

	// test crc
	bufCrc := make([]byte, 4)
	testBuffer.Read(bufCrc)
	c.Assert(binary.LittleEndian.Uint32(bufCrc), DeepEquals, expectedCrc)

	// magic DATA
	magicData := make([]byte, 4)
	testBuffer.Read(magicData)
	c.Assert(magicData, DeepEquals, []byte{'D', 'A', 'T', 'A'})

	// data
	actualData := make([]byte, int32(testLength))
	testBuffer.Read(actualData)
	c.Assert(string(actualData), DeepEquals, testData)

	// extract footer crc32c
	actualFooterCrc := make([]byte, 4)
	testBuffer.Read(actualFooterCrc)
	remainingBytes := testBuffer.Bytes()
	remainingSum, err := crc32c.Crc32c(remainingBytes)
	c.Assert(err, IsNil)
	c.Assert(binary.LittleEndian.Uint32(actualFooterCrc), DeepEquals, remainingSum)

	// sha512
	expectedSha512 := sha512.Sum512([]byte(testData))
	actualSha512 := make([]byte, 64)
	testBuffer.Read(actualSha512)
	c.Assert(actualSha512, DeepEquals, expectedSha512[:])
	log.Println("Length: ", testBuffer.Len())

	// length
	actualLength := make([]byte, 8)
	testBuffer.Read(actualLength)
	c.Assert(testBufferLength, DeepEquals, binary.LittleEndian.Uint64(actualLength))

	// magic INIM
	magicInim := make([]byte, 4)
	testBuffer.Read(magicInim)
	log.Println(magicInim)
	c.Assert(magicInim, DeepEquals, []byte{'I', 'N', 'I', 'M'})

	// ensure no extra data is in the file
	c.Assert(testBuffer.Len(), Equals, 0)
}
