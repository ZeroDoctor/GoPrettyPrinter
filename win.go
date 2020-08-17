// +build windows

package ppt

import "syscall"

func (w *platform) setup(enable bool) error {
	var (
		kernel32Dll    *syscall.LazyDLL  = syscall.NewLazyDLL("Kernel32.dll")
		setConsoleMode *syscall.LazyProc = kernel32Dll.NewProc("SetConsoleMode")
	)

	const ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x4

	var mode uint32
	err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	if err != nil {
		return err
	}

	if enable {
		mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	} else {
		mode &^= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	}

	ret, _, err := setConsoleMode.Call(uintptr(syscall.Stdout), uintptr(mode))
	if ret == 0 {
		return err
	}

	return nil
}
