package config

import (
	"flag"
	"strings"
)

type tagsList []string

func (t *tagsList) String() string {
	return strings.Join(*t, " ")
}

func (t *tagsList) ToStringList() []string {
	return *t
}

func (t *tagsList) Set(value string) error {
	*t = append(*t, value)
	return nil
}

var LogsSenderConfig struct {
	TextLogging    bool
	JournalLogging bool
	Tags           tagsList
	Since          string
}

func ProcessLogsSenderConfigArgs(defaultTextLogging, defaultJournalLogging bool) {
	flag.BoolVar(&LogsSenderConfig.JournalLogging, "with-journal-logging", defaultJournalLogging, "Use journal logging")
	flag.BoolVar(&LogsSenderConfig.TextLogging, "with-text-logging", defaultTextLogging, "Use text logging")
	flag.Var(&LogsSenderConfig.Tags, "tag", "Journalctl tag to filter")
	flag.StringVar(&LogsSenderConfig.Since, "since", "5 hours ago", "Journalctl since flag, same format")
	h := flag.Bool("help", false, "Help message")

	flag.Parse()
	if h != nil && *h {
		printHelpAndExit()
	}
}
