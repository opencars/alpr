package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/opencars/alpr/pkg/api/http"
	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/logger"
	"github.com/opencars/alpr/pkg/recognizer/openalpr"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.toml", "Path to the configuration file")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(configPath)
	if err != nil {
		logger.Fatal(err)
	}

	recognizer, err := openalpr.New(&conf.OpenALPR)
	if err != nil {
		logger.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	addr := ":8080"
	logger.Info("Listening on %s...", addr)
	if err := http.Start(ctx, addr, recognizer); err != nil {
		logger.Fatal(err)
	}
}
