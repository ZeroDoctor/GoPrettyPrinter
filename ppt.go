package ppt

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// Log set custom types
type Log uint16
type Level uint16

// Flags and Colors
const (
	pFATA Log = iota
	pERRO
	pWARN
	pINFO
	pVBSE
	pTRCE
)

const (
	// FUNC enable function call lookup in WhereAmI()
	FUNC = lFlag(1 << iota) // enables function name to be include in log
	// LINE enable line call lookup in WhereAmI()
	LINE // enables line number to be include in log
	// FILE enable line call lookup in WhereAmI()
	FILE // enables file name to be include in log
)

const (
	// lvlFATA the level for logging fatal
	lvlFATA = 1 << iota
	// lvlERRO the level for logging error
	lvlERRO
	// lvlWARN the level for logging warn
	lvlWARN
	// lvlINFO the level for logging info
	lvlINFO
	// lvlVBSE the level for logging verbose
	lvlVBSE
	// lvlTRCE the level for logging trace
	lvlTRCE
)

const (
	FatalLevel   Level = lvlFATA
	ErrorLevel   Level = lvlFATA | lvlERRO
	WarnLevel    Level = lvlFATA | lvlERRO | lvlWARN
	InfoLevel    Level = lvlFATA | lvlERRO | lvlWARN | lvlINFO
	VerboseLevel Level = lvlFATA | lvlERRO | lvlWARN | lvlINFO | lvlVBSE
	TraceLevel   Level = lvlFATA | lvlERRO | lvlWARN | lvlINFO | lvlVBSE | lvlTRCE
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
	displayLog = [6]string{
		"FATL",
		"ERRO",
		"WARN",
		"INFO",
		"VBSE",
		"TRAC",
	}

	LogCallback = [6]func(string, ...interface{}){}

	logType = [6]string{
		FAT + "FATA" + cRST,
		ERR + "ERRO" + cRST,
		WRN + "WARN" + cRST,
		IFO + "INFO" + cRST,
		VER + "VBSE" + cRST,
		TRA + "TRAC" + cRST,
	}

	currentLevel = lvlFATA | lvlERRO | lvlWARN | lvlINFO
	enable       = true
)

// Colors <-
type Colors string

const (
	// Black regular black
	Black = Colors("30")
	// Red regular red
	Red = Colors("31")
	// Green regular green
	Green = Colors("32")
	// Yellow regular yellow
	Yellow = Colors("33")
	// Blue regular blue
	Blue = Colors("34")
	// Magenta regular magenta
	Magenta = Colors("35")
	// Cyan regular cyan
	Cyan = Colors("36")
	// White regular white
	White = Colors("37")

	// Gray a brighter black
	Gray = Colors("90")
	// BRed a brighter red
	BRed = Colors("91")
	// BGreen a brighter green
	BGreen = Colors("92")
	// BYellow a brighter yellow
	BYellow = Colors("93")
	// BBlue a brighter blue
	BBlue = Colors("94")
	// BMagenta a brighter magenta
	BMagenta = Colors("95")
	// BCyan a brighter cyan
	BCyan = Colors("96")
	// BWhite a brighter white
	BWhite = Colors("97")
)

// TogglePrint to console
func TogglePrint() {
	enable = !enable
}

// GetColor return formatted ansi escape color
func GetColor(color Colors) string {
	return begin + string(color) + "m"
}

// ToBackground convert colors to background colors
func ToBackground(color Colors) string {
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

func GetCurrentLevel() int {
	return currentLevel
}

func SetCurrentLevel(lvl Level) {
	currentLevel = int(lvl)
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

// SetFatalColor change the color of fatal log
func SetFatalColor(color Colors) {
	logType[pFATA] = GetColor(color) + "FATA" + cRST
}

// SetErrorColor change the color of error log
func SetErrorColor(color Colors) {
	logType[pERRO] = GetColor(color) + "ERRO" + cRST
}

// SetWarnColor change the color of warn log
func SetWarnColor(color Colors) {
	logType[pWARN] = GetColor(color) + "WARN" + cRST
}

// SetInfoColor change the color of info log
func SetInfoColor(color Colors) {
	logType[pINFO] = GetColor(color) + "INFO" + cRST
}

// SetVerboseColor change the color of verbose log
func SetVerboseColor(color Colors) {
	logType[pVBSE] = GetColor(color) + "VBSE" + cRST
}

// SetTraceColor change the color of trace log
func SetTraceColor(color Colors) {
	logType[pTRCE] = GetColor(color) + "TRAC" + cRST
}

// ###################### Format Log ######################

// Fatalf formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Fatalf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pFATA, str, args)
}

// Errorf formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Errorf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pERRO, str, args)
}

// Warnf formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Warnf(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pWARN, str, args)
}

// Infof formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Infof(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pINFO, str, args)
}

// Verbosef formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Verbosef(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pVBSE, str, args)
}

// Tracef formats according to a format specifier amoung other prefex
// and writes to standard output. It returns the string displayed to console.
func Tracef(msg string, args ...interface{}) string {
	str := fmt.Sprintf(msg, args...)
	return Printer(pTRCE, str, args)
}

// ###################### NewLine Log ######################

// Fatalln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Fatalln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pFATA, msg, args)
}

// Errorln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Errorln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pERRO, msg, args)
}

// Warnln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Warnln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pWARN, msg, args)
}

// Infoln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Infoln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pINFO, msg, args)
}

// Verboseln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Verboseln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pVBSE, msg, args)
}

// Traceln formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the string displayed to console.
func Traceln(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format+"\n", args...)
	return Printer(pTRCE, msg, args)
}

// ###################### Non-Format Log ######################

// Fatal formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Fatal(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pFATA, msg, args)
}

// Error formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Error(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pERRO, msg, args)
}

// Warn formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Warn(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pWARN, msg, args)
}

// Info formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Info(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pINFO, msg, args)
}

// Verbose formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Verbose(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pVBSE, msg, args)
}

// Trace formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the string displayed to console.
func Trace(args ...interface{}) string {
	args = checkPointerType(args...)
	format := getFormatStr(len(args))
	msg := fmt.Sprintf(format, args...)
	return Printer(pTRCE, msg, args)
}

// ###################### The Big Boy on the Block ######################

// Printer output msg to console with a desire log type prefix
func Printer(prefix Log, msg string, args ...interface{}) string {
	if (1 << prefix & currentLevel) == 0 {
		return ""
	}

	result := logType[prefix] + LoggerPrefix() + msg
	if enable {
		fmt.Print(result)
	}
	out := displayLog[prefix] + LoggerPrefix() + msg

	if LogCallback[prefix] != nil {
		LogCallback[prefix](out, args)
	}

	return out
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
