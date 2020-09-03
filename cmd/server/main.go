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
	"github.com/opencars/alpr/pkg/objectstore/minio"
	"github.com/opencars/alpr/pkg/recognizer/openalpr"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.yaml", "Path to the configuration file")

	flag.Parse()

	conf, err := config.New(configPath)
	if err != nil {
		logger.Fatalf("failed read config: %v", err)
	}

	logger.NewLogger(logger.LogLevel(conf.Log.Level), conf.Log.Mode == "dev")

	recognizer, err := openalpr.New(&conf.OpenALPR)
	if err != nil {
		logger.Fatalf("failed to initialize recognizer: %v", err)
	}

	objStore, err := minio.New(&conf.S3)
	if err != nil {
		logger.Fatalf("failed to object store: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-c
		cancel()
	}()

	addr := ":8080"
	logger.Infof("Listening on %s...", addr)
	if err := http.Start(ctx, addr, &conf.Server, recognizer, objStore); err != nil {
		logger.Fatalf("http server failed: %v", err)
	}
}
