package main

import (
	"os"

	"github.com/ori-amizur/introspector/src/commands"
	"github.com/ori-amizur/introspector/src/config"
	"github.com/ori-amizur/introspector/src/util"
)

func main() {
	config.ProcessLogsSenderConfigArgs(false, true)
	util.SetLogging("logs-sender", config.LogsSenderConfig.TextLogging, config.LogsSenderConfig.JournalLogging)
	err := commands.SendLogs(config.LogsSenderConfig.Tags.ToStringList(), config.LogsSenderConfig.Since)
	if err != nil {
		os.Exit(-1)
	}
}
