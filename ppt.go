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
		return ": "
	}
}

// lFlag to prevent others from accidently messing with it
type lFlag uint8

// List of log prefixs
var (
	LoggerFlags    lFlag = 0
	DisplayWarning       = true            // displaying pointer warning is on by default
	Order                = false           // changes the order of LoggerFlags prefix
	LoggerPrefix         = defaultPrefix() // goes after the log type and file|func|line info

	depth      = -1
	displayLog = [5]string{
		"INFO",
		"WARN",
		"ERRO",
		"VBSE",
		"FATL",
	}

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
	_, err := setup(true)
	if err != nil {
		fmt.Println("ERROR (PrettyPrinter): could not setup os")
		panic(err)
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
	logType[pINFO] = GetColor(color) + "INFO" + cRST
}

// SetWarnColor change the color of warn log
func SetWarnColor(color colors) {
	logType[pWARN] = GetColor(color) + "WARN" + cRST
}

// SetErrorColor change the color of error log
func SetErrorColor(color colors) {
	logType[pERROR] = GetColor(color) + "ERROR" + cRST
}

// SetVerboseColor change the color of verbose log
func SetVerboseColor(color colors) {
	logType[pVBOSE] = GetColor(color) + "VBOSE" + cRST
}

// SetFatalColor change the color of fatal log
func SetFatalColor(color colors) {
	logType[pFATAL] = GetColor(color) + "FATAL" + cRST
}

// ###################### Format Log ######################

// Infof formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Infof(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pINFO, str)
}

// Warnf formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Warnf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pWARN, str)
}

// Errorf formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Errorf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pERROR, str)
}

// Verbosef formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Verbosef(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pVBOSE, str)
}

// Fatalf formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Fatalf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pFATAL, str)
}

// ###################### NewLine Log ######################

// Infoln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Infoln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pINFO, msg)
}

// Warnln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Warnln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pWARN, msg)
}

// Errorln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Errorln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pERROR, msg)
}

// Verboseln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Verboseln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pVBOSE, msg)
}

// Fatalln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Fatalln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pFATAL, msg)
}

// ###################### Non-Format Log ######################

// Info formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Info(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pINFO, msg)
}

// Warn formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Warn(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pWARN, msg)
}

// Error formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Error(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pERROR, msg)
}

// Verbose formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Verbose(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pVBOSE, msg)
}

// Fatal formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Fatal(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pFATAL, msg)
}

// ###################### The Big Boy on the Block ######################

// Printer output msg to console with a desire log type prefix
func Printer(prefix uint8, msg string) string {
	result := logType[prefix] + LoggerPrefix() + msg
	fmt.Print(result)
	return displayLog[prefix] + LoggerPrefix() + msg
}

// ###################### Decorator ######################

// lDECOR just another uint8 type
type lDECOR uint8

const (
	afterType lDECOR = iota
	seperator
	afterInfo
)

var (
	decor = [3]string{
		"[",
		"|",
		"]",
	}
)

// Decorator $log_type = [info|warn|error] $log_info = [file|func|line] $extra_prefix is developer define prefix
//
// [0] = after log type i.e. $log_type '[' $log_info $extra_prefix $msg
// [1] = seperator in $log_info i.e. file '|' func '|' line
// [2] = after log info i.e. $log_type $log_info ']' $extra_prefix $msg
func Decorator(args ...string) {
	if len(args) > 3 {
		Warn("PPrinter -- There are only 3 options\n")
	}

	for i, a := range args {
		decor[i] = a
	}
}

// WhereAmI gives the func, file or line of when the log was called
// func, file and line are activiated with LoggerFlags with
// bits FUNC, FILE, LINE
func WhereAmI() string {
	flag := LoggerFlags
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
