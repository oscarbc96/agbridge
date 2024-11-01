package main

import "github.com/oscarbc96/agbridge/pkg/log"

func main() {
	log.Setup(log.LevelDebug)

	log.Info("Starting agbridge")
}
