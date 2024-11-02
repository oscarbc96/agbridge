package main

import (
	"fmt"

	"github.com/oscarbc96/agbridge/pkg/log"
)

func main() {
	options, err := parseFlags()

	log.Setup(options.LogLevel)

	if err != nil {
		log.Fatal(err.Error())
	}

	if options.Version {
		fmt.Println("Application version 1.0.0")
		return
	}

	log.Info("", log.String("resource-id", options.ResourceID), log.String("profile-name", options.ProfileName), log.String("config", options.Config))
}
