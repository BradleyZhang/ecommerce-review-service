package common

import (
	"fmt"
	"log"
	"os"
)

var SysLogger = log.New(os.Stdout, "[SYS] ", log.LstdFlags)

func SysLog(v ...any) {
	SysLogger.Println(v...)
}

var ErrLogger = log.New(os.Stdout, "[ERR] ", log.LstdFlags|log.Llongfile)

func ErrorLog(v ...any) {
	_ = ErrLogger.Output(2, fmt.Sprintln(v...))
}

var DebugLogger = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)

func DebugLog(v ...any) {
	if DebugMode {
		_ = DebugLogger.Output(2, fmt.Sprintln(v...))
	}
}

var FatalLogger = log.New(os.Stdout, "[FATAL] ", log.LstdFlags|log.Llongfile)

func FatalLog(v ...any) {
	_ = FatalLogger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}
