package main

import (
	"fmt"

	"github.com/oscarbc96/agbridge/pkg/log"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	setCustomUsage()
	options, err := parseFlags()

	log.Setup(options.LogLevel)

	if err != nil {
		log.Fatal(err.Error())
	}

	if options.Version {
		fmt.Printf("%s, commit %s, built at %s\n", version, commit, date)

		return
	}

	log.Info("", log.String("resource-id", options.ResourceID), log.String("profile-name", options.ProfileName), log.String("config", options.Config), log.String("listen-address", options.ListenAddress))
}
