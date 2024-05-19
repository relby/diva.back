package logger

import (
	"log"
	"os"
)

var infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var errLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

func Info() *log.Logger {
	return infoLogger
}

func Err() *log.Logger {
	return errLogger
}
