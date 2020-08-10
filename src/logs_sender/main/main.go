package main

import (
	"fmt"
	"os"

	"github.com/ori-amizur/introspector/src/logs_sender"

	"github.com/ori-amizur/introspector/src/config"
	"github.com/ori-amizur/introspector/src/util"
)

func main() {
	config.ProcessLogsSenderConfigArgs(true, true)
	util.SetLogging("logs-sender", config.LogsSenderConfig.TextLogging, config.LogsSenderConfig.JournalLogging)
	err := logs_sender.SendLogs(&logs_sender.LogsSenderExecuter{})
	if err != nil {
		fmt.Println("Failed to run send logs ", err.Error())
		os.Exit(-1)
	}
	fmt.Println("Logs were sent")
}
