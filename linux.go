// +build windows

package ppt

// Linux :
type Linux struct {
	name string
}

func (l *Linux) getName() string {
	return l.name
}

func (l *Linux) setup(enable bool) error {
	return nil
}
