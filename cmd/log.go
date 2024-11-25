package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Level int

// 日志等级
const (
	LevelPanic Level = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

// 超级颜色
const (
	none        = "\033[0m"
	black       = "\033[0;30m"
	darkGray    = "\033[1;30m"
	blue        = "\033[0;34m"
	lightBlue   = "\033[1;34m"
	green       = "\033[0;32m"
	lightGreen  = "\033[1;32m"
	cyan        = "\033[0;36m"
	lightCyan   = "\033[1;36m"
	red         = "\033[0;31m"
	lightRed    = "\033[1;31m"
	purple      = "\033[0;35m"
	lightPurple = "\033[1;35m"
	brown       = "\033[0;33m"
	yellow      = "\033[1;33m"
	lightGray   = "\033[0;37m"
	white       = "\033[1;37m"
)

var (
	Debug  = debugLog.Println
	Debugf = debugLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
	Warn   = warnLog.Println
	Warnf  = warnLog.Printf
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Fatal  = fatalLog.Println
	Fatalf = fatalLog.Printf
	Panic  = panicLog.Println
	Panicf = panicLog.Printf
)

var (
	debugLog = log.New(os.Stdout, fmt.Sprintf("%s[debug]%s", lightGreen, none), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	infoLog  = log.New(os.Stdout, fmt.Sprintf("%s[info ]%s", lightPurple, none), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	warnLog  = log.New(os.Stdout, fmt.Sprintf("%s[warn ]%s", yellow, none), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	errorLog = log.New(os.Stdout, fmt.Sprintf("%s[error]%s", lightRed, none), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	fatalLog = log.New(os.Stdout, fmt.Sprintf("%s[fatal]%s", lightRed, none), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	panicLog = log.New(os.Stdout, fmt.Sprintf("%s[painc]%s", lightRed, none), log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	loggers  = []*log.Logger{debugLog, infoLog, warnLog, errorLog, fatalLog, panicLog}
	mu       sync.Mutex
)

// SetLevel 设置日志等级，打印指定等级以上的日志。默认打印全部
func SetLevel(level Level) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if LevelDebug > level {
		debugLog.SetOutput(io.Discard)
	}
	if LevelInfo > level {
		infoLog.SetOutput(io.Discard)
	}
	if LevelWarn > level {
		warnLog.SetOutput(io.Discard)
	}
	if LevelError > level {
		errorLog.SetOutput(io.Discard)
	}
	if LevelFatal > level {
		fatalLog.SetOutput(io.Discard)
	}
	if LevelPanic > level {
		panicLog.SetOutput(io.Discard)
	}
}
