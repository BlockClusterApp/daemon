package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ClearLogFile() {
	t := time.Now()
	twoDaysBack := t.Add(time.Duration(-48) * time.Hour)

	dateFormat := twoDaysBack.Format("2006-01-02")

	files, err := filepath.Glob(fmt.Sprintf("/tmp/running-logs-%s*", dateFormat))
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			log.Printf("Error deleting log file %s | %s ", f, err.Error())
		}
	}
}
