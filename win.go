// +build windows

package ppt

import (
	"fmt"
	"syscall"
)

func setup(enable bool) (bool, error) {
	fmt.Println("Setting up for windows...")
	var (
		kernel32Dll    *syscall.LazyDLL  = syscall.NewLazyDLL("Kernel32.dll")
		setConsoleMode *syscall.LazyProc = kernel32Dll.NewProc("SetConsoleMode")
	)

	const EnableVirtualTerminalProcessing uint32 = 0x4

	var mode uint32
	err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	if err != nil {
		return false, err
	}

	if enable {
		mode |= EnableVirtualTerminalProcessing
	} else {
		mode &^= EnableVirtualTerminalProcessing
	}

	ret, _, err := setConsoleMode.Call(uintptr(syscall.Stdout), uintptr(mode))
	if ret == 0 {
		return false, err
	}

	return true, nil
}
