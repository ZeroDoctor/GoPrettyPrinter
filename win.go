// +build windows

package ppt

import (
	"fmt"
	"syscall"
)

var (
	begin = "\x1b[1;" // forced 1;
	// IFO default info's log foreground color to regular blue
	IFO = "\x1b[1;36m"
	// WRN default warn's log foreground color to regular yellow
	WRN = "\x1b[1;33m"
	// ERR default error's log foreground color to regular red
	ERR = "\x1b[1;31m"
	// VER default verbose's log foreground color to regular green
	VER = "\x1b[1;32m"
	// FAT default fatal's log color foreground to red and background to yellowish
	FAT       = "\x1b[1;31;103m"
	cRST      = "\x1b[0m"
	clear     = "\x1b[2J"
	rstCursor = "\x1b[1;1H"
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
