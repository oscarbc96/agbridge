package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/oscarbc96/agbridge/pkg/log"
)

func setCustomUsage() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  # Show version
  %[1]s --version

  # Use a specific config file
  %[1]s --config=config.yaml

  # Set profile name with a Rest API ID
  %[1]s --profile-name=myprofile --rest-api-id=12345

  # Set log level to debug
  %[1]s --log-level=debug

  # Set the listen address for the proxy server
  %[1]s --listen-address=:9090

  # Use the default config file (agbridge.yaml or agbridge.yml if they exist)
  %[1]s
`, os.Args[0])
	}
}

type Flags struct {
	Version       bool
	Config        string
	ProfileName   string
	RestAPIID     string
	ListenAddress string
	LogLevel      log.Level
}

func parseFlags() (*Flags, error) {
	setCustomUsage()

	version := flag.Bool("version", false, "Displays the application version and exits.")
	config := flag.String("config", "", "Specifies the path to a configuration file (cannot be used with --profile-name or --rest-api-id).")
	profileName := flag.String("profile-name", "", "Specifies the profile name (requires --rest-api-id to be specified).")
	restAPIID := flag.String("rest-api-id", "", "Specifies the Rest API ID (required if --config is not provided).")
	logLevelStr := flag.String("log-level", "info", "Sets the log verbosity level. Options: debug, info, warn, error, fatal.")
	listenAddress := flag.String("listen-address", ":8080", "Address where the proxy server will listen for incoming requests.")

	flag.Parse()

	if *version {
		return &Flags{Version: true}, nil
	}

	logLevel, err := log.ParseLogLevel(*logLevelStr)
	if err != nil {
		return nil, err
	}

	flags := &Flags{
		Version:       *version,
		Config:        *config,
		ProfileName:   *profileName,
		RestAPIID:     *restAPIID,
		ListenAddress: *listenAddress,
		LogLevel:      logLevel,
	}

	if _, _, err := net.SplitHostPort(*listenAddress); err != nil {
		return flags, fmt.Errorf("invalid listen address format")
	}

	// Check if a custom config file is specified and verify its existence
	if *config != "" {
		// If config is specified, it must not be combined with other flags
		if *profileName != "" || *restAPIID != "" {
			return flags, errors.New("--config cannot be combined with --profile-name or --rest-api-id")
		}

		// Ensure the config file exists
		if _, err := os.Stat(*config); os.IsNotExist(err) {
			return flags, errors.New("config file does not exist")
		}
	} else {
		// If config is not specified, check the necessity of Rest API ID
		if *restAPIID == "" && *profileName != "" {
			return flags, errors.New("--profile-name requires --rest-api-id to be specified")
		}

		// If no config and no Rest API ID, check for default config files
		if *restAPIID == "" {
			if _, err := os.Stat("agbridge.yaml"); os.IsNotExist(err) {
				if _, err := os.Stat("agbridge.yml"); os.IsNotExist(err) {
					return flags, errors.New("please provide --rest-api-id, --config, or ensure agbridge.yaml or agbridge.yml exists")
				} else {
					flags.Config = "agbridge.yml" // Default to agbridge.yml if it exists
				}
			} else {
				flags.Config = "agbridge.yaml" // Default to agbridge.yaml if it exists
			}
		}
	}

	return flags, nil
}
