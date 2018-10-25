package helpers

import (
	"io"
	"log"
	"os"
	"sync"
)

type Logger struct {
	filename string
	*log.Logger
}

var logger *Logger
var once sync.Once

func GetLogger() *Logger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func createLogger() *Logger {
	file, _ := os.OpenFile("/tmp/running-logs.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	f := io.MultiWriter(file, os.Stdout)

	return &Logger{
		Logger: log.New(f, "B-Agent ", log.LUTC),
	}
}
