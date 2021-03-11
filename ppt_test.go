package ppt

import (
	"testing"
	"time"
)

// TestFormat <-
func TestColor(t *testing.T) {
	Infoln("Sup")
	SetInfoColor(Blue)
	Infoln("Sup now in blue")
}

func TestPrefix(t *testing.T) {
	Infoln("Nothing here")
	LoggerPrefix = func() string {
		return " [" + time.Now().Format("15:04:05") + "]"
	}
	Infoln("Something here")
}

func TestWhereAmI(t *testing.T) {
	Infoln("Not here")
	LoggerPrefix = func() string {
		return " [" + time.Now().Format("15:04:05") + "]" + WhereAmI() + ": "
	}
	Infoln("Still not here")
	LoggerFlags = FILE | LINE
	Infoln("Here I am")
}
