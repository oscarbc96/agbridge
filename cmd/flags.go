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
	Config        string
	ListenAddress string
	LogLevel      log.Level
	ProfileName   string
	Region        string
	RestAPIID     string
	Version       bool
}

func parseFlags() (*Flags, error) {
	setCustomUsage()

	version := flag.Bool("version", false, "Displays the application version and exits.")
	config := flag.String("config", "", "Specifies the path to a configuration file (cannot be used with --profile-name, --rest-api-id, or --region).")
	profileName := flag.String("profile-name", "", "Specifies the profile name (requires --rest-api-id and --region to be specified).")
	restAPIID := flag.String("rest-api-id", "", "Specifies the Rest API ID (required if --config is not provided).")
	region := flag.String("region", "", "Specifies the AWS region to use with --profile-name and --rest-api-id.")
	logLevelStr := flag.String("log-level", "info", "Sets the log verbosity level. Options: debug, info, warn, error, fatal.")
	listenAddress := flag.String("listen-address", ":8080", "Address where the proxy server will listen for incoming requests.")

	flag.Parse()

	if *version {
		return &Flags{Version: true}, nil
	}

	logLevel, err := log.ParseLogLevel(*logLevelStr)
	if err != nil {
		return &Flags{LogLevel: logLevel}, err
	}

	flags := &Flags{
		Version:       *version,
		Config:        *config,
		ProfileName:   *profileName,
		RestAPIID:     *restAPIID,
		ListenAddress: *listenAddress,
		LogLevel:      logLevel,
		Region:        *region,
	}

	// Validate listen address format
	if _, _, err := net.SplitHostPort(*listenAddress); err != nil {
		return flags, fmt.Errorf("invalid listen address format: %w", err)
	}

	// Check if a custom config file is specified and verify its existence
	if *config != "" {
		if *profileName != "" || *restAPIID != "" || *region != "" {
			return flags, errors.New("--config cannot be combined with --profile-name, --rest-api-id, or --region")
		}

		if _, err := os.Stat(*config); os.IsNotExist(err) {
			return flags, fmt.Errorf("config file does not exist: %w", err)
		}
	} else {
		// If no --config, check the rules for --rest-api-id, --region, and --profile-name

		// --profile-name requires both --region and --rest-api-id
		if *profileName != "" && (*restAPIID == "" || *region == "") {
			return flags, errors.New("--profile-name requires both --region and --rest-api-id to be specified")
		}

		// --region requires --rest-api-id
		if *region != "" && *restAPIID == "" {
			return flags, errors.New("--region requires --rest-api-id to be specified")
		}

		// If neither --config nor --rest-api-id is provided, fallback to default config file check
		if *restAPIID == "" {
			if _, err := os.Stat("agbridge.yaml"); os.IsNotExist(err) {
				if _, err := os.Stat("agbridge.yml"); os.IsNotExist(err) {
					return flags, errors.New("please provide --rest-api-id, --config, or ensure agbridge.yaml or agbridge.yml exists")
				} else {
					flags.Config = "agbridge.yml"
				}
			} else {
				flags.Config = "agbridge.yaml"
			}
		}
	}

	return flags, nil
}
