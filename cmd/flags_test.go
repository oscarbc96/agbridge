package main

import (
	"flag"
	"os"
	"testing"

	"github.com/oscarbc96/agbridge/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		expErr  string
		expOpts *Flags
	}{
		{
			name:   "Version only",
			args:   []string{"cmd", "--version"},
			expErr: "",
			expOpts: &Flags{
				Version: true,
			},
		},
		{
			name:   "No Flags with default config",
			args:   []string{"cmd"},
			expErr: "",
			expOpts: &Flags{
				Config:        "agbridge.yaml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid RestAPIID only",
			args:   []string{"cmd", "--rest-api-id", "12345"},
			expErr: "",
			expOpts: &Flags{
				RestAPIID:     "12345",
				ProfileName:   "",
				Region:        "",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid Region and RestAPIID",
			args:   []string{"cmd", "--region", "eu-west-1", "--rest-api-id", "12345"},
			expErr: "",
			expOpts: &Flags{
				RestAPIID:     "12345",
				ProfileName:   "",
				Region:        "eu-west-1",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid Region, RestAPIID, and ProfileName",
			args:   []string{"cmd", "--region", "eu-west-1", "--rest-api-id", "12345", "--profile-name", "patata"},
			expErr: "",
			expOpts: &Flags{
				RestAPIID:     "12345",
				ProfileName:   "patata",
				Region:        "eu-west-1",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "ProfileName without Region and RestAPIID",
			args:   []string{"cmd", "--profile-name", "patata"},
			expErr: "`--profile-name` requires both `--region` and `--rest-api-id` to be specified",
			expOpts: &Flags{
				ProfileName:   "patata",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Region without RestAPIID",
			args:   []string{"cmd", "--region", "eu-west-1"},
			expErr: "`--region` requires `--rest-api-id` to be specified",
			expOpts: &Flags{
				Region:        "eu-west-1",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Config and RestAPIID",
			args:   []string{"cmd", "--config", "config.yaml", "--rest-api-id", "12345"},
			expErr: "`--config` cannot be combined with `--profile-name`, `--rest-api-id`, or `--region`",
			expOpts: &Flags{
				Config:        "config.yaml",
				RestAPIID:     "12345",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Config and ProfileName",
			args:   []string{"cmd", "--config", "config.yaml", "--profile-name", "testprofile"},
			expErr: "`--config` cannot be combined with `--profile-name`, `--rest-api-id`, or `--region`",
			expOpts: &Flags{
				Config:        "config.yaml",
				ProfileName:   "testprofile",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Config and Region",
			args:   []string{"cmd", "--config", "config.yaml", "--region", "eu-west-1"},
			expErr: "`--config` cannot be combined with `--profile-name`, `--rest-api-id`, or `--region`",
			expOpts: &Flags{
				Config:        "config.yaml",
				Region:        "eu-west-1",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Invalid Config File",
			args:   []string{"cmd", "--config", "nonexistent.yaml"},
			expErr: "config file does not exist: stat nonexistent.yaml: no such file or directory",
			expOpts: &Flags{
				Config:        "nonexistent.yaml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "No Default Config Files",
			args:   []string{"cmd"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Only agbridge.yml exists",
			args:   []string{"cmd"},
			expErr: "",
			expOpts: &Flags{
				Config:        "agbridge.yml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Only agbridge.yaml exists",
			args:   []string{"cmd"},
			expErr: "",
			expOpts: &Flags{
				Config:        "agbridge.yaml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid Listen Address",
			args:   []string{"cmd", "--listen-address", ":9090"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":9090",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Invalid Listen Address",
			args:   []string{"cmd", "--listen-address", "qwerty"},
			expErr: "invalid listen address format: address qwerty: missing port in address",
			expOpts: &Flags{
				ListenAddress: "qwerty",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid LogLevel - Debug",
			args:   []string{"cmd", "--log-level", "debug"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelDebug,
			},
		},
		{
			name:   "Valid LogLevel - Info",
			args:   []string{"cmd", "--log-level", "info"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid LogLevel - Warn",
			args:   []string{"cmd", "--log-level", "warn"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelWarn,
			},
		},
		{
			name:   "Valid LogLevel - Error",
			args:   []string{"cmd", "--log-level", "error"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelError,
			},
		},
		{
			name:   "Valid LogLevel - Fatal",
			args:   []string{"cmd", "--log-level", "fatal"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelFatal,
			},
		},
		{
			name:   "Invalid LogLevel",
			args:   []string{"cmd", "--log-level", "verbose"},
			expErr: "invalid log level: must be one of debug, info, warn, error, fatal",
			expOpts: &Flags{
				LogLevel: log.LevelInfo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()

			if tt.name == "No Flags with default config" || tt.name == "Only agbridge.yaml exists" {
				file, err := os.Create("agbridge.yaml")
				require.NoError(t, err)
				file.Close()
				defer os.Remove("agbridge.yaml")
			}

			if tt.name == "Only agbridge.yml exists" {
				os.Remove("agbridge.yaml")
				file, err := os.Create("agbridge.yml")
				require.NoError(t, err)
				file.Close()
				defer os.Remove("agbridge.yml")
			}

			os.Args = tt.args

			opts, err := parseFlags()

			if tt.expErr == "" {
				require.NoError(t, err, "Expected no error")
			} else {
				require.Error(t, err, "Expected error")
				require.EqualError(t, err, tt.expErr)
			}
			assert.Equal(t, tt.expOpts, opts, "Options mismatch")
		})
	}
}
