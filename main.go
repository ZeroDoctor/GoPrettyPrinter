package ppt

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// Flags and Colors
const (
	pINFO = iota
	pWARN
	pERROR
	pVBOSE
	pFATAL

	FUNC = lFlag(1 << (iota - 5)) // enables function name to be include in log
	LINE                          // enables line number to be include in log
	FILE                          // enables file name to be include in log
)

func defaultPrefix() func() string {
	return func() string {
		return " "
	}
}

// lFlag to prevent others from accidently messing with it
type lFlag uint8

// some Log prefixs
var (
	begin     = "\033[1;"        // forced 1;
	IFO       = "\033[1;34m"     // default info's log foreground color to regular blue
	WRN       = "\033[1;33m"     // default warn's log foreground color to regular yellow
	ERR       = "\033[1;31m"     // default error's log foreground color to regular red
	VER       = "\033[1;32m"     // default verbose's log foreground color to regular green
	FAT       = "\033[1;31;103m" // default fatal's log color foreground to red and background to yellowish
	cRST      = "\033[0m"
	clear     = "\033[2J"
	rstCursor = "\033[1;1H"

	LoggerFlags    lFlag = 0
	DisplayWarning       = true            // displaying pointer warning is on by default
	Order                = false           // changes the order of LoggerFlags prefix
	LoggerPrefix         = defaultPrefix() // goes after the log type and file|func|line info

	depth   = -1
	logType = [5]string{
		IFO + "INFO" + cRST,
		WRN + "WARN" + cRST,
		ERR + "ERROR" + cRST,
		VER + "VBOSE" + cRST,
		FAT + "FATAL" + cRST,
	}
)

type colors string

const (
	// Black regular black
	Black = colors("30")
	// Red regular red
	Red = colors("31")
	// Green regular green
	Green = colors("32")
	// Yellow regular yellow
	Yellow = colors("33")
	// Blue regular blue
	Blue = colors("34")
	// Magenta regular magenta
	Magenta = colors("35")
	// Cyan regular cyan
	Cyan = colors("36")
	// White regular white
	White = colors("37")

	// Gray a brighter black
	Gray = colors("90")
	// BRed a brighter red
	BRed = colors("91")
	// BGreen a brighter green
	BGreen = colors("92")
	// BYellow a brighter yellow
	BYellow = colors("93")
	// BBlue a brighter blue
	BBlue = colors("94")
	// BMagenta a brighter magenta
	BMagenta = colors("95")
	// BCyan a brighter cyan
	BCyan = colors("96")
	// BWhite a brighter white
	BWhite = colors("97")
)

// GetColor return formatted ansi escape color
func GetColor(color colors) string {
	return begin + string(color) + "m"
}

// ToBackground convert colors to background colors
func ToBackground(color colors) string {
	colorNum, err := strconv.Atoi(string(color))
	if err != nil {
		fmt.Println("ERROR (PrettyPrinter): could not convert color to background color")
		return string(color)
	}

	colorNum += 10
	newColor := strconv.Itoa(colorNum)
	return newColor
}

// Init setup colors based on OS
func Init() {
	windows, err := setup(true)
	if err != nil {
		fmt.Println("ERROR (PrettyPrinter): could not setup os")
		panic(err)
	}

	if windows {
		begin = "\x1b[1;"
		IFO = "\x1b[1;34m"
		WRN = "\x1b[1;33m"
		ERR = "\x1b[1;31m"
		VER = "\x1b[1;32m"
		FAT = "\x1b[1;31;103m"
		cRST = "\x1b[0m"
		clear = "\x1b[2J"
		rstCursor = "\x1b[1;1H"

		logType = [5]string{
			IFO + "INFO" + cRST,
			WRN + "WARN" + cRST,
			ERR + "ERROR" + cRST,
			VER + "VBOSE" + cRST,
			FAT + "FATAL" + cRST,
		}
	}
}

// Clear clears console
func Clear() {
	fmt.Print(clear)
}

// ResetCursor move cursor to the upper left corner
func ResetCursor() {
	fmt.Print(rstCursor)
}

// ResetColor changes consoles' colors back to normal
func ResetColor() string {
	return cRST
}

// SetInfoColor change the color of info log
func SetInfoColor(color colors) {
	logType[pINFO] = string(color) + "INFO" + cRST
}

// SetWarnColor change the color of warn log
func SetWarnColor(color colors) {
	logType[pWARN] = string(color) + "WARN" + cRST
}

// SetErrorColor change the color of error log
func SetErrorColor(color colors) {
	logType[pERROR] = string(color) + "ERROR" + cRST
}

// SetVerboseColor change the color of verbose log
func SetVerboseColor(color colors) {
	logType[pVBOSE] = string(color) + "VBOSE" + cRST
}

// SetFatalColor change the color of fatal log
func SetFatalColor(color colors) {
	logType[pFATAL] = string(color) + "FATAL" + cRST
}

// ###################### Format Log ######################

func Infof(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pINFO, str)
}

func Warnf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pWARN, str)
}

func Errorf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pERROR, str)
}

func Verbosef(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pVBOSE, str)
}

func Fatalf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pFATAL, str)
}

// ###################### NewLine Log ######################

func Infoln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pINFO, msg)
}

func Warnln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pWARN, msg)
}

func Errorln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pERROR, msg)
}

func Verboseln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pVBOSE, msg)
}

func Fatalln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pFATAL, msg)
}

// ###################### Non-Format Log ######################

func Info(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pINFO, msg)
}

func Warn(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pWARN, msg)
}

func Error(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pERROR, msg)
}

func Verbose(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pVBOSE, msg)
}

func Fatal(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pFATAL, msg)
}

// ###################### The Big Boy on the Block ######################

// Printer output msg to console with a desire log type prefix
func Printer(prefix uint8, msg string) string {
	result := logType[prefix] + swap(Order) + msg
	fmt.Print(result)
	return result
}

// ###################### Decorator ######################

// lDECOR just another uint8 type
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
		" [",
		"]-",
		"|",
	}
)

// Decorator $log_type = [info|warn|error] $log_info = [file|func|line] $extra_prefix is developer define prefix
//
// [0] = between prefix i.e. $log_type $log_info ':' $extra_prefix $msg
// [1] = after log type i.e. $log_type '{' $log_info $extra_prefix $msg
// [2] = after log info i.e. $log_type $log_info '}' $extra_prefix $msg
// [3] = seperator in $log_info i.e. file '|' func '|' line
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
