package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/opencars/seedwork/logger"

	"github.com/opencars/alpr/pkg/api/http"
	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/domain/service"
	"github.com/opencars/alpr/pkg/queue/nats"
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

	pub, err := nats.New(conf.EventAPI.Address(), conf.EventAPI.Enabled)
	if err != nil {
		logger.Fatalf("nats: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-c
		cancel()
	}()

	svc := service.NewCustomerService(recognizer, pub)

	addr := ":8080"
	logger.Infof("Listening on %s...", addr)
	if err := http.Start(ctx, addr, &conf.Server, svc); err != nil {
		logger.Fatalf("http server failed: %v", err)
	}
}
