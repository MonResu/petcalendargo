package logger

import (
	"fmt"
	"log"
	"os"
)

var (
    infoLogger *log.Logger
    errorLogger *log.Logger
	logFile *os.File
)

func LogInit(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("LogInit: %w", err)
	}
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logFile = file
	return nil
}

func PrintInfo(info string) {
    if infoLogger != nil {
        infoLogger.Output(2, info)
    } else {
        log.Println("INFO:", info) 
    }
}

func PrintError(errorMsg string) {
    if errorLogger != nil {
        errorLogger.Output(2, errorMsg)
    } else {
        log.Println("ERROR:", errorMsg) 
    }
}

func Close() error {
    if logFile != nil {
        return logFile.Close()
    }
    return nil
}