package main

import (
	"testing"

	"github.com/oscarbc96/agbridge/pkg/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		expErr  string
		expOpts *Flags
		setup   func(t *testing.T, fs afero.Fs)
	}{
		{
			name:   "Version only",
			args:   []string{"--version"},
			expErr: "",
			expOpts: &Flags{
				Version: true,
			},
		},
		{
			name:   "No Flags with default config",
			args:   []string{},
			expErr: "",
			expOpts: &Flags{
				Config:        "agbridge.yaml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
			setup: func(t *testing.T, fs afero.Fs) {
				require.NoError(t, afero.WriteFile(fs, DefaultConfigFileYaml, []byte("dummy"), 0o644))
			},
		},
		{
			name:   "Valid RestAPIID only",
			args:   []string{"--rest-api-id", "12345"},
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
			args:   []string{"--region", "eu-west-1", "--rest-api-id", "12345"},
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
			args:   []string{"--region", "eu-west-1", "--rest-api-id", "12345", "--profile-name", "patata"},
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
			args:   []string{"--profile-name", "patata"},
			expErr: "`--profile-name` requires both `--region` and `--rest-api-id` to be specified",
			expOpts: &Flags{
				ProfileName:   "patata",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Region without RestAPIID",
			args:   []string{"--region", "eu-west-1"},
			expErr: "`--region` requires `--rest-api-id` to be specified",
			expOpts: &Flags{
				Region:        "eu-west-1",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Config and RestAPIID",
			args:   []string{"--config", "config.yaml", "--rest-api-id", "12345"},
			expErr: "`--config` cannot be combined with `--profile-name`, `--rest-api-id`, `--region`, or `--stage-name`",
			expOpts: &Flags{
				Config:        "config.yaml",
				RestAPIID:     "12345",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Config and ProfileName",
			args:   []string{"--config", "config.yaml", "--profile-name", "testprofile"},
			expErr: "`--config` cannot be combined with `--profile-name`, `--rest-api-id`, `--region`, or `--stage-name`",
			expOpts: &Flags{
				Config:        "config.yaml",
				ProfileName:   "testprofile",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Config and Region",
			args:   []string{"--config", "config.yaml", "--region", "eu-west-1"},
			expErr: "`--config` cannot be combined with `--profile-name`, `--rest-api-id`, `--region`, or `--stage-name`",
			expOpts: &Flags{
				Config:        "config.yaml",
				Region:        "eu-west-1",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Config and StageName",
			args:   []string{"--config", "config.yaml", "--stage-name", "test"},
			expErr: "`--config` cannot be combined with `--profile-name`, `--rest-api-id`, `--region`, or `--stage-name`",
			expOpts: &Flags{
				Config:        "config.yaml",
				StageName:     "test",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Invalid Config File",
			args:   []string{"--config", "nonexistent.yaml"},
			expErr: "config file does not exist: open nonexistent.yaml: file does not exist",
			expOpts: &Flags{
				Config:        "nonexistent.yaml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "No Default Config Files",
			args:   []string{},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Only agbridge.yml exists",
			args:   []string{},
			expErr: "",
			expOpts: &Flags{
				Config:        "agbridge.yml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
			setup: func(t *testing.T, fs afero.Fs) {
				require.NoError(t, afero.WriteFile(fs, DefaultConfigFileYml, []byte("dummy"), 0o644))
			},
		},
		{
			name:   "Only agbridge.yaml exists",
			args:   []string{},
			expErr: "",
			expOpts: &Flags{
				Config:        "agbridge.yaml",
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
			setup: func(t *testing.T, fs afero.Fs) {
				require.NoError(t, afero.WriteFile(fs, DefaultConfigFileYaml, []byte("dummy"), 0o644))
			},
		},
		{
			name:   "Valid Listen Address",
			args:   []string{"--listen-address", ":9090"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":9090",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Invalid Listen Address",
			args:   []string{"--listen-address", "qwerty"},
			expErr: "invalid listen address format: address qwerty: missing port in address",
			expOpts: &Flags{
				ListenAddress: "qwerty",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid LogLevel - Debug",
			args:   []string{"--log-level", "debug"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelDebug,
			},
		},
		{
			name:   "Valid LogLevel - Info",
			args:   []string{"--log-level", "info"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelInfo,
			},
		},
		{
			name:   "Valid LogLevel - Warn",
			args:   []string{"--log-level", "warn"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelWarn,
			},
		},
		{
			name:   "Valid LogLevel - Error",
			args:   []string{"--log-level", "error"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelError,
			},
		},
		{
			name:   "Valid LogLevel - Fatal",
			args:   []string{"--log-level", "fatal"},
			expErr: "please provide `--rest-api-id`, `--config`, or ensure agbridge.yaml or agbridge.yml exists",
			expOpts: &Flags{
				ListenAddress: ":8080",
				LogLevel:      log.LevelFatal,
			},
		},
		{
			name:   "Invalid LogLevel",
			args:   []string{"--log-level", "verbose"},
			expErr: "invalid log level: must be one of debug, info, warn, error, fatal",
			expOpts: &Flags{
				LogLevel: log.LevelInfo,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := afero.NewMemMapFs()
			if tt.setup != nil {
				tt.setup(t, fs)
			}

			opts, err := parseFlags(fs, tt.args)

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
