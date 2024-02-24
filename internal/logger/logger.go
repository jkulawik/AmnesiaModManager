package logger

import (
	"log"
	"os"
	"runtime"
)

const colorReset = "\033[0m"

var (
	logFile *os.File
	Warn    *log.Logger
	Error   *log.Logger
	Info    *log.Logger
)

func init() {
	// defer logFile.Close() // FIXME this might be needed for file logging
	var (
		infoM = "[INFO]    "
		warnM = "[WARNING] "
		errM  = "[ERROR]   "
	)

	if runtime.GOOS == "linux" {
		infoM = "\033[36m" + infoM + colorReset
		warnM = "\033[33m" + warnM + colorReset
		errM = "\033[31m" + errM + colorReset
	}

	writer := os.Stderr
	Info = log.New(writer, infoM, log.Lshortfile)
	Warn = log.New(writer, warnM, log.Lshortfile)
	Error = log.New(writer, errM, log.Lshortfile)

	if runtime.GOOS == "windows" {
		newFlags := log.Ldate | log.Ltime | log.Lshortfile
		Info.SetFlags(newFlags)
		Warn.SetFlags(newFlags)
		Error.SetFlags(newFlags)

		var err error
		logFile, err = os.OpenFile("modmanager.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			Error.Printf("error opening log file: %v", err)
		}

		Info.SetOutput(logFile)
		Warn.SetOutput(logFile)
		Error.SetOutput(logFile)

	}
}
