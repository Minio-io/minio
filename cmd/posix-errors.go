/*
 * Minio Cloud Storage, (C) 2016, 2017 Minio, Inc.
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
	"os"
	"runtime"
	"syscall"
)

// Function not implemented error
func isSysErrNoSys(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err == syscall.ENOSYS {
			return true
		}
	}
	return false
}

// Not supported error
func isSysErrOpNotSupported(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err == syscall.EOPNOTSUPP {
			return true
		}
	}
	return false
}

// No space left on device error
func isSysErrNoSpace(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err == syscall.ENOSPC {
			return true
		}
	}
	return false
}

// Input/output error
func isSysErrIO(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err == syscall.EIO {
			return true
		}
	}
	return false
}

// Check if the given error corresponds to EISDIR (is a directory).
func isSysErrIsDir(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err == syscall.EISDIR {
			return true
		}
	}
	return false
}

// Check if the given error corresponds to ENOTDIR (is not a directory).
func isSysErrNotDir(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err == syscall.ENOTDIR {
			return true
		}
	}
	return false
}

// Check if the given error corresponds to the ENAMETOOLONG (name too long).
func isSysErrTooLong(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Err == syscall.ENAMETOOLONG {
			return true
		}
	}
	return false
}

// Check if the given error corresponds to ENOTEMPTY for unix
// and ERROR_DIR_NOT_EMPTY for windows (directory not empty).
func isSysErrNotEmpty(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		if runtime.GOOS == globalWindowsOSName {
			if errno, _ok := pathErr.Err.(syscall.Errno); _ok && errno == 0x91 {
				// ERROR_DIR_NOT_EMPTY
				return true
			}
		}
		if pathErr.Err == syscall.ENOTEMPTY {
			return true
		}
	}
	return false
}

// Check if the given error corresponds to the specific ERROR_PATH_NOT_FOUND for windows
func isSysErrPathNotFound(err error) bool {
	if runtime.GOOS != globalWindowsOSName {
		return false
	}
	if pathErr, ok := err.(*os.PathError); ok {
		if errno, _ok := pathErr.Err.(syscall.Errno); _ok && errno == syscall.ESRCH {
			// ERROR_PATH_NOT_FOUND
			return true
		}
	}
	return false
}

// Check if the given error corresponds to the specific ERROR_INVALID_HANDLE for windows
func isSysErrHandleInvalid(err error) bool {
	if runtime.GOOS != globalWindowsOSName {
		return false
	}
	// Check if err contains ERROR_INVALID_HANDLE errno
	errno, ok := err.(syscall.Errno)
	return ok && errno == syscall.ENXIO
}

func isSysErrCrossDevice(err error) bool {
	e, ok := err.(*os.LinkError)
	return ok && e.Err == syscall.EXDEV
}
