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

func TestParseFlags_VersionOnly(t *testing.T) {
	resetFlags()

	os.Args = []string{"cmd", "--version"}

	opts, err := parseFlags()

	require.NoError(t, err, "Expected no error")
	assert.True(t, opts.Version, "Expected version flag to be true")
}

func TestParseFlags_IncompatibleFlagsConfigAndProfile(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--config", "config.yaml", "--profile-name", "testprofile"}

	_, err := parseFlags()

	require.Error(t, err, "Expected incompatible flags error")
	assert.EqualError(t, err, "--config cannot be combined with --profile-name or --rest-api-id")
}

func TestParseFlags_IncompatibleFlagsConfigAndRestAPIID(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--config", "config.yaml", "--rest-api-id", "12345"}

	_, err := parseFlags()

	require.Error(t, err, "Expected incompatible flags error")
	assert.EqualError(t, err, "--config cannot be combined with --profile-name or --rest-api-id")
}

func TestParseFlags_IncompatibleFlagsConfigRestAPIIDAndProfile(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--config", "config.yaml", "--rest-api-id", "12345", "--profile-name", "testprofile"}

	_, err := parseFlags()

	require.Error(t, err, "Expected incompatible flags error")
	assert.EqualError(t, err, "--config cannot be combined with --profile-name or --rest-api-id")
}

func TestParseFlags_IncompatibleFlagsProfile(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--profile-name", "testprofile"}

	_, err := parseFlags()

	require.Error(t, err, "Expected incompatible flags error")
	assert.EqualError(t, err, "--profile-name requires --rest-api-id to be specified")
}

func TestParseFlags_ValidRestAPIID(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--rest-api-id", "12345"}

	opts, err := parseFlags()

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "12345", opts.RestAPIID, "Rest API ID mismatch")
	assert.Empty(t, opts.ProfileName, "Expected profile name to be empty")
	assert.Equal(t, log.LevelInfo, opts.LogLevel, "Default log level mismatch")
	assert.Equal(t, ":8080", opts.ListenAddress, "Default listen address mismatch")
}

func TestParseFlags_ValidRestAPIIDAndProfileName(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--rest-api-id", "12345", "--profile-name", "testprofile"}

	opts, err := parseFlags()

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, "12345", opts.RestAPIID, "Rest API ID mismatch")
	assert.Equal(t, "testprofile", opts.ProfileName, "Profile name mismatch")
	assert.Equal(t, log.LevelInfo, opts.LogLevel, "Default log level mismatch")
	assert.Equal(t, ":8080", opts.ListenAddress, "Default listen address mismatch")
}

func TestParseFlags_NoFlags(t *testing.T) {
	resetFlags()

	// Create a temporary default config file
	file, err := os.Create("agbridge.yaml")
	require.NoError(t, err)
	file.Close()
	defer os.Remove("agbridge.yaml")

	os.Args = []string{"cmd"}

	opts, err := parseFlags()

	require.NoError(t, err, "Expected no error")
	assert.Equal(t, log.LevelInfo, opts.LogLevel, "Default log level mismatch")
	assert.Equal(t, "agbridge.yaml", opts.Config, "Config file mismatch")
	assert.Equal(t, ":8080", opts.ListenAddress, "Default listen address mismatch")
}

func TestParseFlags_ValidConfig(t *testing.T) {
	resetFlags()

	// Create a temporary config file to simulate a valid --config file
	tmpFile, err := os.CreateTemp("", "config.yaml")
	require.NoError(t, err, "Failed to create temporary config file")
	defer os.Remove(tmpFile.Name()) // Clean up after test

	os.Args = []string{"cmd", "--config", tmpFile.Name()}

	opts, err := parseFlags()

	require.NoError(t, err, "Expected no error")
	assert.EqualValues(t, log.LevelInfo, opts.LogLevel, "Default log level mismatch")
	assert.Equal(t, tmpFile.Name(), opts.Config, "Config file path mismatch")
	assert.Equal(t, ":8080", opts.ListenAddress, "Default listen address mismatch")
}

func TestParseFlags_InvalidConfigFileNotExist(t *testing.T) {
	resetFlags()

	os.Args = []string{"cmd", "--config", "nonexistent.yaml"}

	_, err := parseFlags()

	require.Error(t, err, "Expected error for nonexistent config file")
	assert.EqualError(t, err, "config file does not exist")
}

func TestParseFlags_NoDefaultConfigFiles(t *testing.T) {
	resetFlags()

	os.Args = []string{"cmd"}

	_, err := parseFlags()

	require.Error(t, err, "Expected error when no default config files are present")
	assert.EqualError(t, err, "please provide --rest-api-id, --config, or ensure agbridge.yaml or agbridge.yml exists")
}

func TestParseFlags_OnlyAgbridgeYmlExists(t *testing.T) {
	resetFlags()

	// Ensure agbridge.yaml does not exist and create agbridge.yml
	os.Remove("agbridge.yaml")
	tmpFile, err := os.Create("agbridge.yml")
	require.NoError(t, err, "Failed to create agbridge.yml")
	defer os.Remove("agbridge.yml")
	tmpFile.Close()

	// Run with no flags to trigger the default config file check
	os.Args = []string{"cmd"}

	opts, err := parseFlags()

	require.NoError(t, err, "Expected no error when agbridge.yml exists")
	assert.Equal(t, "agbridge.yml", opts.Config, "Config file path should default to agbridge.yml")
	assert.Equal(t, log.LevelInfo, opts.LogLevel, "Default log level mismatch")
	assert.Equal(t, ":8080", opts.ListenAddress, "Default listen address mismatch")
}

func TestParseFlags_ValidLogLevel(t *testing.T) {
	resetFlags()

	// Create a temporary default config file to avoid missing file error
	file, err := os.Create("agbridge.yaml")
	require.NoError(t, err)
	file.Close()
	defer os.Remove("agbridge.yaml")

	os.Args = []string{"cmd", "--log-level", "debug"}

	opts, err := parseFlags()

	require.NoError(t, err, "Expected no error")
	assert.EqualValues(t, log.LevelDebug, opts.LogLevel, "Default log level mismatch")
	assert.Equal(t, ":8080", opts.ListenAddress, "Default listen address mismatch")
}

func TestParseFlags_InvalidLogLevel(t *testing.T) {
	resetFlags()

	os.Args = []string{"cmd", "--log-level", "verbose"}

	_, err := parseFlags()

	require.Error(t, err, "Expected invalid log level error")
	assert.EqualError(t, err, "invalid log level: must be one of debug, info, warn, error, fatal")
}

func TestParseFlags_ValidListenAddress(t *testing.T) {
	resetFlags()

	// Create a temporary default config file
	file, err := os.Create("agbridge.yaml")
	require.NoError(t, err)
	file.Close()
	defer os.Remove("agbridge.yaml")

	os.Args = []string{"cmd", "--listen-address", ":9090"}

	opts, err := parseFlags()
	require.NoError(t, err, "Expected no error")
	assert.Equal(t, ":9090", opts.ListenAddress, "Listen address mismatch")
	assert.Equal(t, log.LevelInfo, opts.LogLevel, "Default log level mismatch")
}

func TestParseFlags_InvalidListenAddress(t *testing.T) {
	resetFlags()

	os.Args = []string{"cmd", "--listen-address", "qwerty"}

	_, err := parseFlags()

	require.Error(t, err, "Expected invalid listen address error")
	assert.EqualError(t, err, "invalid listen address format")
}

func TestParseFlags_DefaultListenAddress(t *testing.T) {
	resetFlags()

	// Create a temporary default config file to avoid missing file error
	file, err := os.Create("agbridge.yaml")
	require.NoError(t, err)
	file.Close()
	defer os.Remove("agbridge.yaml")

	os.Args = []string{"cmd"}

	opts, err := parseFlags()
	require.NoError(t, err, "Expected no error")
	assert.Equal(t, ":8080", opts.ListenAddress, "Default listen address mismatch")
	assert.Equal(t, log.LevelInfo, opts.LogLevel, "Default log level mismatch")
}
