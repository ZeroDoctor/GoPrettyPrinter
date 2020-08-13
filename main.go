package ppt

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Flags and Colors
const (
	pINFO = iota
	pWARN
	pERROR
	pVBOSE
	pFATAL

	FUNC = lFlag(1 << iota)
	LINE
	FILE
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
	BLU  = "\033[1;34m"
	YEL  = "\033[1;33m"
	RED  = "\033[1;31m"
	GRE  = "\033[1;32m"
	FAT  = "\033[1;31;103m"
	cRST = "\033[0m"

	LoggerFlags    lFlag = 0
	DisplayWarning       = true
	Order                = false
	LoggerPrefix         = defaultPrefix() // goes after the log type and file|func|line info

	depth   = -1
	logType = [5]string{
		BLU + "INFO" + cRST,
		YEL + "WARN" + cRST,
		RED + "ERROR" + cRST,
		GRE + "VBOSE" + cRST,
		FAT + "FATAL" + cRST,
	}
)

// SetInfoColor :
func SetInfoColor(color string) {
	logType[pINFO] = color + "INFO" + cRST
}

// SetWarnColor :
func SetWarnColor(color string) {
	logType[pWARN] = color + "INFO" + cRST
}

// SetErrorColor :
func SetErrorColor(color string) {
	logType[pERROR] = color + "INFO" + cRST
}

// SetVerboseColor :
func SetVerboseColor(color string) {
	logType[pVBOSE] = color + "INFO" + cRST
}

// SetFatalColor :
func SetFatalColor(color string) {
	logType[pFATAL] = color + "INFO" + cRST
}

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
	Printer(pINFO, msg)
}

// Warnln :
func Warnln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(pWARN, msg)
}

// Errorln :
func Errorln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(pERROR, msg)
}

// Verboseln :
func Verboseln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(pVBOSE, msg)
}

// Fatalln :
func Fatalln(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	Printer(pFATAL, msg)
}

// ###################### Non-Format Log ######################

// Info :
func Info(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(pINFO, msg)
}

// Warn :
func Warn(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(pWARN, msg)
}

// Error :
func Error(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(pERROR, msg)
}

// Verbose :
func Verbose(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(pVBOSE, msg)
}

// Fatal :
func Fatal(args ...interface{}) {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	Printer(pFATAL, msg)
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
		fmt.Print(logType[pWARN] + ": PPrinter -- Couldn't recover [function/file/line]\n")
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
		funcStr := runtime.FuncForPC(function).Name()
		start := strings.LastIndex(funcStr, "/") // remove folder path
		format += funcStr[start+1:]
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
