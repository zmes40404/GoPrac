package util

import (
	"log"
	"os"
)

var (
	INFO *log.Logger
	WARN *log.Logger
	ERROR *log.Logger
	DEBUG *log.Logger
)

func init() {
	loggfile, err := os.OpenFile(".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if  err != nil {
		panic(err)
	}
	INFO = log.New(loggfile, "INFO: ", log.LstdFlags|log.Llongfile)
	WARN = log.New(loggfile, "WARN: ", log.LstdFlags|log.Llongfile)
	ERROR = log.New(loggfile, "ERROR: ", log.LstdFlags|log.Llongfile)
	DEBUG = log.New(loggfile, "DEBUG: ", log.LstdFlags|log.Llongfile)
}