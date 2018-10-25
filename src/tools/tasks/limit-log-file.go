package tasks

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"log"
	"os"
	"os/exec"
)

func ClearLogFile() {
	logger := helpers.GetLogger()
	logger.Println("G:TASK Clearing logs")

	if _, err := os.Stat("/tmp/blockcluster.log"); !os.IsNotExist(err) {
		// path/to/whatever exists
		exec.Command("touch", "blockcluster.log")
		if err != nil {
			log.Fatalf("Error creating file failed with %s\n", err)
		}
	}

	cmd := exec.Command("tail", "-n 10000", "/tmp/running-logs.log")
	log.Println(cmd.Output())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
