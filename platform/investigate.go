package platform

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/hpcloud/tail"
)

func StartInvestigate(extractedDirectory string, ok1 chan bool) {
	ch1 := make(chan bool)

	// setup log file
	file, err := os.Create("investigate.log")
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("bin/investigate")
	cmd.Dir = extractedDirectory + "/siren-investigate"
	cmd.Stdout = file
	//cmd.Stderr = os.Stderr

	go checkError(ch1)

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	<-ch1
	ok1 <- true
}

func checkError(ch1 chan bool) {
	// we initialize our InvestigateLog json
	var investigateLog InvestigateLog
	// we initialize elasticsearch_client error count
	es_client_error_count := 0
	// we initialize max error count
	max_error_count := 30

	file, err := tail.TailFile("investigate2.log", tail.Config{
		Follow: true,
		ReOpen: true,
		Poll:   true, // With poll = true appended the problem does NOT occur, so suspect it has something to do with fsnotify
	})
	if err != nil {
		fmt.Println(err)
	}

	for line := range file.Lines {

		// we unmarshal our byteArray which contains our
		// jsonFile's content into 'investigateLog' which we defined above
		json.Unmarshal([]byte(line.Text), &investigateLog)

		tags_length := len(investigateLog.Tags)
		tag := investigateLog.Tags[0]
		nextTag := investigateLog.Tags[tags_length-1]
		// index is the index where we are
		// element is the element from someSlice for where we are
		if strings.EqualFold(tag, "error") || strings.EqualFold(tag, "warning") {
			if strings.EqualFold(nextTag, "elasticsearch") {
				if !(es_client_error_count <= max_error_count) {
					fmt.Println(strings.Join(investigateLog.Tags, " ") + investigateLog.Message)
					log.Fatal(strings.Join(investigateLog.Tags, " ") + investigateLog.Message)
					syscall.Kill(investigateLog.Pid, syscall.SIGKILL)
					os.Exit(1)
				}
				es_client_error_count++
			} else {
				fmt.Println(strings.Join(investigateLog.Tags, " ") + investigateLog.Message)
				log.Fatal(strings.Join(investigateLog.Tags, " ") + investigateLog.Message)
				syscall.Kill(investigateLog.Pid, syscall.SIGKILL)
				os.Exit(1)
			}
		}
		if strings.EqualFold(tag, "info") && strings.EqualFold(nextTag, "console") {
			if strings.EqualFold(investigateLog.Message, "Template [template:kibi-html-angular] successfully loaded") {
				fmt.Println("No error or warning in investgate console")
				syscall.Kill(investigateLog.Pid, syscall.SIGKILL)
				ch1 <- true
			}
		}
	}
}

// InvestigateLog struct which contains a tags, message
type InvestigateLog struct {
	Type      string   `json:"type"`
	Timestamp string   `json:"@timestamp"`
	Tags      []string `json:"tags"`
	Pid       int      `json:"pid"`
	State     string   `json:"state"`
	Message   string   `json:"message"`
	PrevState string   `json:"prevState"`
	PrevMsg   string   `json:"prevMsg"`
}
