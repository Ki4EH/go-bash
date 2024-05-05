package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

func init() {
	output, err := OutPutFile()
	if err != nil {
		fmt.Errorf("error creating logfile %s", err)
	}

	infoLog = log.New(output, "INFO: ", log.Ldate|log.Ltime)
	errorLog = log.New(output, "ERROR: ", log.Ldate|log.Ltime)

}

// OutPutFile Создаю директорию для логов и лог файл на текущий день
func OutPutFile() (*os.File, error) {
	dir, _ := os.Getwd()
	logDir := filepath.Join(dir, "logs")

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, os.ModeDir)
		if err != nil {
			return nil, fmt.Errorf("error creating dir %s", err)

		}
	}

	today := time.Now().Format("2006-01-02")
	filePath := filepath.Join(logDir, today+".log")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error creating log file %s", err)
	}
	return file, nil
}

func Info(v ...interface{}) {
	infoLog.Printf(" %s", fmt.Sprintln(v...))
}

func Error(v ...interface{}) {
	file, line := getFileAndLine(2)
	errorLog.Printf("[%s:%d] %s", file, line, fmt.Sprintln(v...))
}

func Fatal(v ...interface{}) {
	file, line := getFileAndLine(2)
	errorLog.Printf("[%s:%d] %s", file, line, fmt.Sprintln(v...))
	os.Exit(1)
}

func getFileAndLine(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "??"
		line = 0
	}
	return filepath.Base(file), line
}
