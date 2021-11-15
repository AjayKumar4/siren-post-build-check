package platform

import (
	"fmt"
	"os"
	"os/exec"
)

func StartFederate(extractedDirectory string, ok2 chan bool) {
	// setup log file
	file, err := os.Create("elasticsearch.log")
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("bin/elasticsearch")
	cmd.Dir = extractedDirectory + "/elasticsearch"
	cmd.Stdout = file
	//cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	ok2 <- true
}

// FederateLog struct which contains a tags, message
type FederateLog struct {
	Timestamp []string
	Type      []string
	Status    []string
	Cluster   []string
	Message   string
}
