package ppt

import (
	"log"
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

func TestAllType(t *testing.T) {
	// Init()

	Fatalln("This is fatal")
	Errorln("This is a error")
	Warnln("This is a warn")
	Infoln("This is a info")
	Verboseln("This is verbose")
	Traceln("This is a trace")
}

func TestLogLevel(t *testing.T) {
	SetCurrentLevel(FatalLevel)
	Fatalln("This is fatal")
	Errorln("This is a error")
	Warnln("This is a warn")
	Infoln("This is a info")
	Verboseln("This is verbose")
	Traceln("This is a trace")
	log.Println("only fatal")

	SetCurrentLevel(InfoLevel)
	Fatalln("This is fatal")
	Errorln("This is a error")
	Warnln("This is a warn")
	Infoln("This is a info")
	Verboseln("This is verbose")
	Traceln("This is a trace")
	log.Println("only info level and up")

	SetCurrentLevel(TraceLevel)
	Fatalln("This is fatal")
	Errorln("This is a error")
	Warnln("This is a warn")
	Infoln("This is a info")
	Verboseln("This is verbose")
	Traceln("This is a trace")
	log.Println("only trace level and up")
}
