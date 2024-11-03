package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/oscarbc96/agbridge/pkg/log"
	"github.com/oscarbc96/agbridge/pkg/proxy"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func loadProxyConfig(flags *Flags) (*proxy.Config, error) {
	if flags.RestAPIID != "" {
		return proxy.NewConfig(flags.RestAPIID, flags.ProfileName, flags.Region), nil
	}

	cfg, err := proxy.LoadConfig(flags.Config)
	if err != nil {
		return nil, fmt.Errorf("couldn't load config file: %w", err)
	}

	return cfg, nil
}

func main() {
	flags, err := parseFlags()
	// Setup logging, before raising errors of flags parsing
	log.Setup(flags.LogLevel)
	if err != nil {
		log.Fatal(err.Error())
	}

	if flags.Version {
		fmt.Printf("%s, commit %s, built at %s\n", version, commit, date)
		return
	}

	cfg, err := loadProxyConfig(flags)
	if err != nil {
		log.Fatal("Failed to load configuration", log.Err(err))
	}

	handlerMapping, err := cfg.Validate()
	if err != nil {
		log.Fatal("Couldn't validate config", log.Err(err))
	}

	proxy := proxy.NewProxy(flags.ListenAddress, handlerMapping)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("Starting proxy", log.String("addr", proxy.Addr()))
		if err := proxy.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Proxy server error", log.Err(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutdown signal received, stopping proxy server...")
	if err := proxy.Shutdown(); err != nil {
		log.Fatal("Failed to stop proxy server gracefully", log.Err(err))
	}
}
