package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	iter := func(imagesPath, name string) error {
		f, err := os.Open(imagesPath + "/" + name)
		if err != nil {
			logger.Errorf("failed to open: %s", err)
			return nil
		}
		defer f.Close()

		result, err := w.r.Recognize(f)
		if err != nil {
			logger.Errorf("error: %s", err)
			return nil
		}

		if len(result) == 0 {
			logger.Errorf("result not found")
			return nil
		}

		var failed bool
		for _, number := range result {
			logger.Infof("detected: %s", number.Plate)

			err = w.db.Recognition().Create(ctx, &model.Recognition{
				ImageKey: "plates/" + name,
				Plate:    number.Plate,
			})

			if err != nil {
				logger.Errorf("db error: %s", err)
				failed = true
			}
		}

		if !failed {
			logger.Infof("success: %s", name)

			if err := os.Remove(imagesPath + "/" + name); err != nil {
				logger.Errorf("db error: %s", err)
				return nil
			}
		}

		return nil
	}

	entries, err := os.ReadDir(imagesPath)
	if err != nil {
		return err
	}

	logger.Infof("count of files in the folder: %d", len(entries))

	for _, e := range entries {
		if err := iter(imagesPath, e.Name()); err != nil {
			logger.Errorf("err: %s", err)
			return nil
		}
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

	logger.Infof("parsing images")
	w := newWorker(sub, obj, sql, recognizer)
	if err := w.process(ctx, imagesPath); err != nil {
		logger.Fatalf("process: %v", err)
	}
}
