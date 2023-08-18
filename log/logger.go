package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	DISABLE
)

var (
	errLog   = log.New(os.Stdout, "\033[31m[ERROR]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[INFO]\033[0m ", log.LstdFlags|log.Lshortfile)
	debugLog = log.New(os.Stdout, "\033[33m[DEBUG]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errLog, infoLog, debugLog}
	mu       sync.Mutex

	Debug  = debugLog.Println
	Debugf = debugLog.Printf
	Error  = errLog.Println
	Errorf = errLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

func SetLevel(lv LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	// Reset all logger first
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if DEBUG < lv {
		debugLog.SetOutput(ioutil.Discard)
	}
	if INFO < lv {
		infoLog.SetOutput(ioutil.Discard)
	}

	if ERROR < lv {
		errLog.SetOutput(ioutil.Discard)
	}

}
