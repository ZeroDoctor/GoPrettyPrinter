// +build windows

package ppt

import "syscall"

// Linux :
type Linux struct {
	name string
}

func (l *Linux) getName() string {
	return l.name
}

func (l *Linux) setup(stream syscall.Handle, enable bool) error {
	return nil
}
