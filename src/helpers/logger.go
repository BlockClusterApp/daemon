package helpers

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	filename string
	*log.Logger
}

var logger *Logger
var once sync.Once

func GetLogger() *Logger {
	gocron.Every(1).Hour().Do(func() {
		logger = createLogger()
	})
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func createLogger() *Logger {
	t := time.Now()
	timeDisplay := t.Format("2006-01-02-15")

	filePath := fmt.Sprintf("/tmp/running-logs-%s.log", timeDisplay)

	file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	f := io.MultiWriter(file, os.Stdout)

	return &Logger{
		Logger: log.New(f, "B-Agent ", log.LUTC),
	}
}
