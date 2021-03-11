// +build linux darwin

package ppt

var (
	begin = "\033[1;" // forced 1;
	// IFO default info's log foreground color to regular blue
	IFO = "\033[1;36m"
	// WRN default warn's log foreground color to regular yellow
	WRN = "\033[1;33m"
	// ERR default error's log foreground color to regular red
	ERR = "\033[1;31m"
	// VER default verbose's log foreground color to regular green
	VER = "\033[1;32m"
	// FAT default fatal's log color foreground to red and background to yellowish
	FAT       = "\033[1;31;103m"
	cRST      = "\033[0m"
	clear     = "\033[2J"
	rstCursor = "\033[1;1H"
)

func setup(enable bool) (bool, error) {
	return false, nil
}
