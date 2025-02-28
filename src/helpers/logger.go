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

	once.Do(func() {
		gocron.Every(1).Hour().Do(func() {
			logger = createLogger()
			log.Printf("Created logger")
		})
		logger = createLogger()
		log.Printf("Created logger")
	})
	return logger
}

func createLogger() *Logger {
	t := time.Now()
	bc := GetBlockclusterInstance()

	timeDisplay := t.Format("2006-01-02-15")

	filePath := fmt.Sprintf("/tmp/running-logs-%s.log", timeDisplay)

	file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	var f io.Writer

	if os.Getenv("SHOW_LOGS") != "" {
		f = io.MultiWriter(file, os.Stdout)
	} else {
		f = io.MultiWriter(file)
	}


	return &Logger{
		Logger: log.New(f, fmt.Sprintf("B-Agent %s ", bc.Metadata.ClientID), log.LUTC),
	}
}

func RefreshLogger() {
	logger = createLogger()
}