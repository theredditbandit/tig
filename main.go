package main

import (
	"tig/cmd"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetReportCaller(true)
	log.SetReportTimestamp(false)
	log.SetFormatter(log.TextFormatter)
	log.SetLevel(log.DebugLevel)

	cmd.Execute()
}
