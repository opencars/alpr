package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opencars/seedwork/logger"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/domain"
	"github.com/opencars/alpr/pkg/domain/model"
	"github.com/opencars/alpr/pkg/objectstore/minio"
	"github.com/opencars/alpr/pkg/queue/nats"
	"github.com/opencars/alpr/pkg/recognizer/openalpr"
	"github.com/opencars/alpr/pkg/store/sqlstore"
)

type worker struct {
	sub domain.Subscriber
	obj domain.ObjectStore
	db  domain.Store
	r   domain.Recognizer

	http *http.Client
}

func (w *worker) process(ctx context.Context, imagesPath string) error {
	entries, err := os.ReadDir(imagesPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		f, err := os.Open(imagesPath + "/" + e.Name())
		if err != nil {
			return err
		}

		result, err := w.r.Recognize(f)
		if err != nil {
			return err
		}

		time.Sleep(time.Second)
		fmt.Println(result)
	}

	return nil
}

func newWorker(sub domain.Subscriber, obj domain.ObjectStore, db domain.Store, recognizer domain.Recognizer) *worker {
	return &worker{
		sub: sub,
		obj: obj,
		db:  db,
		r:   recognizer,
		http: &http.Client{
			Timeout: model.ClientTimeOut,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: model.ClientTimeOut,
				}).DialContext,
			},
		},
	}
}

func main() {
	var configPath string
	var imagesPath string

	flag.StringVar(&configPath, "config", "./config/config.yaml", "Path to the configuration file")
	flag.StringVar(&imagesPath, "images", "./plates/", "Images")

	flag.Parse()

	conf, err := config.New(configPath)
	if err != nil {
		logger.Fatalf("config: %v", err)
	}

	logger.NewLogger(logger.LogLevel(conf.Log.Level), conf.Log.Mode == "dev")

	recognizer, err := openalpr.New(&conf.OpenALPR)
	if err != nil {
		logger.Fatalf("failed to initialize recognizer: %v", err)
	}

	obj, err := minio.New(&conf.S3)
	if err != nil {
		logger.Fatalf("minio: %v", err)
	}

	sql, err := sqlstore.New(&conf.DB)
	if err != nil {
		logger.Fatalf("store: %v", err)
	}

	sub, err := nats.New(conf.NATS.Address(), conf.NATS.Enabled)
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

	logger.Infof("listening events from queue")
	w := newWorker(sub, obj, sql, recognizer)
	if err := w.process(ctx, imagesPath); err != nil {
		logger.Fatalf("process: %v", err)
	}
}
