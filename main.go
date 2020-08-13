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

	FUNC = lFlag(1 << iota)
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

// lFlag : to prevent others from accidently messing with it
type lFlag uint8

// some Log prefixs
var (
	LoggerFlags    lFlag = 0
	DisplayWarning       = true
	Order                = false
	LoggerPrefix         = defaultPrefix() // goes after the log type and file|func|line info

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
	depth = 5
	str := fmt.Sprintf(msg, args...)
	Info(str)
	depth = -1
}

// Warnf :
func Warnf(msg string, args ...interface{}) {
	depth = 5
	str := fmt.Sprintf(msg, args...)
	Warn(str)
	depth = -1
}

// Errorf :
func Errorf(msg string, args ...interface{}) {
	depth = 5
	str := fmt.Sprintf(msg, args...)
	Error(str)
	depth = -1
}

// Verbosef :
func Verbosef(msg string, args ...interface{}) {
	depth = 5
	str := fmt.Sprintf(msg, args...)
	Verbose(str)
	depth = -1
}

// Fatalf :
func Fatalf(msg string, args ...interface{}) {
	depth = 5
	str := fmt.Sprintf(msg, args...)
	Fatal(str)
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
	fmt.Print(logType[prefix] + swap(Order) + msg)
}

// ###################### Decorator ######################

// lDECOR : just another uint8 type
type lDECOR uint8

const (
	betweenPrefix lDECOR = iota
	afterType
	afterInfo
	seperator
)

var (
	decor = [4]string{
		":",
		" {",
		"}->",
		"|",
	}
)

// Decorator : $log_type = [info|warn|error] $log_info = [file|func|line] $extra_prefix is developer define prefix
//	[0] = between prefix i.e. $log_type $log_info ':' $extra_prefix $msg
//	[1] = after log type i.e. $log_type '{' $log_info $extra_prefix $msg
//	[2] = after log info i.e. $log_type $log_info '}' $extra_prefix $msg
//	[3] = seperator in $log_info i.e. file '|' func '|' line
func Decorator(args ...string) {
	if len(args) > 5 {
		Warn("PPrinter -- There are only 4 options\n")
	}

	for i, a := range args {
		decor[i] = a
	}
}

// ###################### Utils ######################

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

func swap(order bool) string {
	if order {
		return LoggerPrefix() + decor[betweenPrefix] + whereAmI(LoggerFlags)
	}

	return whereAmI(LoggerFlags) + decor[betweenPrefix] + LoggerPrefix()
}

func whereAmI(flag lFlag) string {
	if flag == 0 {
		return ""
	}

	if depth == -1 {
		depth = 4
	}

	function, file, line, ok := runtime.Caller(depth)
	if DisplayWarning && !ok {
		fmt.Print(logType[WARN] + ": PPrinter -- Couldn't recover [function/file/line]\n")
		return ""
	}

	format := decor[afterType]
	if flag&(FILE) != 0 {
		format += fileOnly(file)
		if flag&(FUNC|LINE) != 0 {
			format += decor[seperator]
		}
	}

	if flag&(FUNC) != 0 {
		format += runtime.FuncForPC(function).Name()
		if flag&(LINE) != 0 {
			format += decor[seperator]
		}
	}

	if flag&(LINE) != 0 {
		format += fmt.Sprintf("%d", line)
	}

	return format + decor[afterInfo]
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
