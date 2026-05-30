package native

/*
#cgo CFLAGS: -I${SRCDIR}
#include <stdlib.h>
#include "native.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func CreateSnapshot(volume, name string) (string, error) {
	cVolume := C.CString(volume)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cVolume))
	defer C.free(unsafe.Pointer(cName))

	return call(func(buffer *C.char, length C.int) C.int {
		return C.nas_recovery_create_snapshot(cVolume, cName, buffer, length)
	})
}

func RestoreFile(snapshotID, sourcePath, targetPath string) (string, error) {
	cSnapshot := C.CString(snapshotID)
	cSource := C.CString(sourcePath)
	cTarget := C.CString(targetPath)
	defer C.free(unsafe.Pointer(cSnapshot))
	defer C.free(unsafe.Pointer(cSource))
	defer C.free(unsafe.Pointer(cTarget))

	return call(func(buffer *C.char, length C.int) C.int {
		return C.nas_recovery_restore_file(cSnapshot, cSource, cTarget, buffer, length)
	})
}

func RollbackVolume(snapshotID, volume string) (string, error) {
	cSnapshot := C.CString(snapshotID)
	cVolume := C.CString(volume)
	defer C.free(unsafe.Pointer(cSnapshot))
	defer C.free(unsafe.Pointer(cVolume))

	return call(func(buffer *C.char, length C.int) C.int {
		return C.nas_recovery_rollback_volume(cSnapshot, cVolume, buffer, length)
	})
}

func ApplyNetworkProfile(nicID, mode, ipv4, gateway string, simulateFailure bool) (string, error) {
	cNic := C.CString(nicID)
	cMode := C.CString(mode)
	cIPv4 := C.CString(ipv4)
	cGateway := C.CString(gateway)
	defer C.free(unsafe.Pointer(cNic))
	defer C.free(unsafe.Pointer(cMode))
	defer C.free(unsafe.Pointer(cIPv4))
	defer C.free(unsafe.Pointer(cGateway))

	var failure C.int
	if simulateFailure {
		failure = 1
	}

	return call(func(buffer *C.char, length C.int) C.int {
		return C.nas_network_apply_profile(cNic, cMode, cIPv4, cGateway, failure, buffer, length)
	})
}

func CreateVM(name string, cpu, memoryGB, diskGB int, networkID string) (string, error) {
	cName := C.CString(name)
	cNetworkID := C.CString(networkID)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cNetworkID))

	return call(func(buffer *C.char, length C.int) C.int {
		return C.nas_vm_create(cName, C.int(cpu), C.int(memoryGB), C.int(diskGB), cNetworkID, buffer, length)
	})
}

func ChangeVMPower(vmID, action string) (string, error) {
	cVMID := C.CString(vmID)
	cAction := C.CString(action)
	defer C.free(unsafe.Pointer(cVMID))
	defer C.free(unsafe.Pointer(cAction))

	return call(func(buffer *C.char, length C.int) C.int {
		return C.nas_vm_power(cVMID, cAction, buffer, length)
	})
}

func call(fn func(*C.char, C.int) C.int) (string, error) {
	buffer := make([]byte, 256)
	code := fn((*C.char)(unsafe.Pointer(&buffer[0])), C.int(len(buffer)))
	message := C.GoString((*C.char)(unsafe.Pointer(&buffer[0])))
	if code != 0 {
		return message, fmt.Errorf(message)
	}
	return message, nil
}
