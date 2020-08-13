package pprinter

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Flags and Colors
const (
	INFO = iota
	WARN
	ERROR
	VBOSE
	FATAL

	FUNC = uint8(1 << iota)
	LINE
	FILE

	BLU = "\033[1;34m"
	YEL = "\033[1;33m"
	RED = "\033[1;31m"
	GRE = "\033[1;32m"
	FAT = "\033[1;31;103m"
	RST = "\033[0m"
)

func defaultPrefix() func() string {
	return func() string {
		return " "
	}
}

// some Log prefixs
var (
	LoggerFlags    uint8 = 0
	DisplayWarning       = true
	LoggerPrefix         = defaultPrefix()

	depth   = -1
	logType = [5]string{
		BLU + "INFO" + RST,
		YEL + "WARN" + RST,
		RED + "ERROR" + RST,
		GRE + "VBOSE" + RST,
		FAT + "FATAL" + RST,
	}
)

// ###################### Format Log ######################

// Infof : logs info
func Infof(msg string, args ...interface{}) {
	depth = 4
	str := fmt.Sprintf(msg, args...)
	Info(str)
	depth = -1
}

// Warnf :
func Warnf(msg string, args ...interface{}) {
	depth = 4
	str := fmt.Sprintf(msg, args...)
	Warn(str)
	depth = -1
}

// Errorf :
func Errorf(msg string, args ...interface{}) {
	depth = 4
	str := fmt.Sprintf(msg, args...)
	Error(str)
	depth = -1
}

// Verbosef :
func Verbosef(msg string, args ...interface{}) {
	depth = 4
	str := fmt.Sprintf(msg, args...)
	Verbose(str)
	depth = -1
}

// ###################### NewLine Log ######################

// Infoln :
func Infoln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(INFO, msg)
}

// Warnln :
func Warnln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(WARN, msg)
}

// Errorln :
func Errorln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(ERROR, msg)
}

// Verboseln :
func Verboseln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(VBOSE, msg)
}

// Fatalln :
func Fatalln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(FATAL, msg)
}

// ###################### Non-Format Log ######################

// Info :
func Info(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(INFO, msg)
}

// Warn :
func Warn(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(WARN, msg)
}

// Error :
func Error(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(ERROR, msg)
}

// Verbose :
func Verbose(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(VBOSE, msg)
}

// Fatal :
func Fatal(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(FATAL, msg)
}

// ###################### The Big Boy on the Block ######################

// Printer :
func Printer(prefix uint8, msg string) {
	fmt.Print(logType[prefix] + whereAmI(LoggerFlags) + ":" + LoggerPrefix() + msg)
}

func checkPointerType(args ...interface{}) []interface{} {

	if !DisplayWarning {
		return args
	}

	for _, a := range args {
		switch reflect.ValueOf(a).Kind() {
		case reflect.Ptr:
			Warn("PPrinter -- is Displaying a Pointer\n")
		}
	}

	return args
}

func getFormatStr(length int) string {

	format := ""
	for i := 0; i < length; i++ {
		format += "%+v "
	}

	return format[:len(format)-1]
}

func fileOnly(str string) string {
	i := strings.LastIndex(str, "/")
	if i == -1 {
		return str
	}

	return str[i+1:]
}

func whereAmI(flag uint8) string {
	if flag == 0 {
		return ""
	}

	if depth == -1 {
		depth = 3
	}

	function, file, line, ok := runtime.Caller(depth)
	if DisplayWarning && !ok {
		fmt.Print(logType[WARN] + " PPrinter -- Couldn't recover [function/file/line] \n")
		return ""
	}

	format := ">-{"
	if flag&(FILE) != 0 {
		format += fileOnly(file)
		if flag&(FUNC|LINE) != 0 {
			format += "|"
		}
	}

	if flag&(FUNC) != 0 {
		format += runtime.FuncForPC(function).Name()
		if flag&(LINE) != 0 {
			format += "|"
		}
	}

	if flag&(LINE) != 0 {
		format += fmt.Sprintf("%d", line)
	}

	return format + "}-<"
}
