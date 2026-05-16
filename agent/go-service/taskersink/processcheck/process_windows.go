//go:build windows

package processcheck

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	TH32CS_SNAPPROCESS = 0x00000002
	MAX_PATH           = 260
)

// PROCESSENTRY32 is the Windows structure for process enumeration
type PROCESSENTRY32 struct {
	Size              uint32
	Usage             uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	Threads           uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

var (
	kernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procCreateToolhelp32  = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First    = kernel32.NewProc("Process32FirstW")
	procProcess32Next     = kernel32.NewProc("Process32NextW")
)

// enumerateProcessNames returns a list of all running process names
func enumerateProcessNames() ([]string, error) {
	snapshot, _, err := procCreateToolhelp32.Call(TH32CS_SNAPPROCESS, 0)
	if snapshot == 0 || snapshot == uintptr(windows.InvalidHandle) {
		return nil, err
	}
	defer windows.CloseHandle(windows.Handle(snapshot))

	var pe32 PROCESSENTRY32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	ret, _, err := procProcess32First.Call(snapshot, uintptr(unsafe.Pointer(&pe32)))
	if ret == 0 {
		return nil, err
	}

	var names []string
	for {
		name := windows.UTF16ToString(pe32.ExeFile[:])
		if name != "" {
			names = append(names, name)
		}

		ret, _, _ = procProcess32Next.Call(snapshot, uintptr(unsafe.Pointer(&pe32)))
		if ret == 0 {
			break
		}
	}

	return names, nil
}
