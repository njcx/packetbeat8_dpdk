// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// MACHINE GENERATED BY 'go generate' COMMAND; DO NOT EDIT

package procs

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modiphlpapi = windows.NewLazySystemDLL("iphlpapi.dll")

	procGetExtendedTcpTable = modiphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdpTable = modiphlpapi.NewProc("GetExtendedUdpTable")
)

func _GetExtendedTcpTable(pTcpTable uintptr, pdwSize *uint32, bOrder bool, ulAf uint32, tableClass uint32, reserved uint32) (code syscall.Errno, err error) {
	var _p0 uint32
	if bOrder {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := syscall.Syscall6(procGetExtendedTcpTable.Addr(), 6, uintptr(pTcpTable), uintptr(unsafe.Pointer(pdwSize)), uintptr(_p0), uintptr(ulAf), uintptr(tableClass), uintptr(reserved))
	code = syscall.Errno(r0)
	if code == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetExtendedUdpTable(pTcpTable uintptr, pdwSize *uint32, bOrder bool, ulAf uint32, tableClass uint32, reserved uint32) (code syscall.Errno, err error) {
	var _p0 uint32
	if bOrder {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := syscall.Syscall6(procGetExtendedUdpTable.Addr(), 6, uintptr(pTcpTable), uintptr(unsafe.Pointer(pdwSize)), uintptr(_p0), uintptr(ulAf), uintptr(tableClass), uintptr(reserved))
	code = syscall.Errno(r0)
	if code == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
