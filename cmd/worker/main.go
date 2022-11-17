package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opencars/seedwork/logger"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/domain/model"
	"github.com/opencars/alpr/pkg/objectstore"
	"github.com/opencars/alpr/pkg/objectstore/minio"
	"github.com/opencars/alpr/pkg/queue"
	"github.com/opencars/alpr/pkg/queue/nats"
	"github.com/opencars/alpr/pkg/store"
	"github.com/opencars/alpr/pkg/store/sqlstore"
)

const (
	// MaxImageSize equals to 5 MB.
	MaxImageSize = 5 << 20

	// ClientTimeOut equals to 10 seconds.
	ClientTimeOut = 10 * time.Second
)

type worker struct {
	sub queue.Subscriber
	obj objectstore.ObjectStore
	db  store.Store

	http *http.Client
}

func (w *worker) process(ctx context.Context) error {
	events, err := w.sub.Subscribe(ctx)
	if err != nil {
		return err
	}

	iter := func() (bool, error) {
		event, ok := <-events
		if !ok {
			return false, nil
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, event.URL, nil)
		if err != nil {
			return true, err
		}

		resp, err := w.http.Do(req)
		if err != nil {
			return true, err
		}

		bodyWithLimit := io.LimitReader(resp.Body, MaxImageSize+1)

		var buff bytes.Buffer
		if _, err = io.CopyN(&buff, bodyWithLimit, bytes.MinRead); err != nil {
			return true, err
		}

		typ := http.DetectContentType(buff.Bytes())
		if typ != "image/jpeg" {
			return true, fmt.Errorf("unsupported content-type: %s", typ)
		}

		_, err = buff.ReadFrom(bodyWithLimit)
		if err != nil {
			return true, err
		}

		if buff.Len() > MaxImageSize {
			return true, fmt.Errorf("image is too big")
		}

		reader := bytes.NewReader(buff.Bytes())
		hash := md5.New()

		if _, err := io.Copy(hash, reader); err != nil {
			return true, err
		}

		_, err = reader.Seek(0, 0)
		if err != nil {
			return true, err
		}

		key := fmt.Sprintf("plates/%s.jpeg", hex.EncodeToString(hash.Sum(nil)))

		err = w.obj.Put(ctx, key, reader)
		if err == nil {
			err = w.db.Recognition().Create(ctx, &model.Recognition{
				ImageKey: key,
				Plate:    event.Number,
			})
		}

		return true, err
	}

	for {
		resume, err := iter()
		if !resume {
			return err
		} else if err != nil {
			logger.Errorf("iter: %v", err)
		}
	}
}

func newWorker(sub queue.Subscriber, obj objectstore.ObjectStore, db store.Store) *worker {
	return &worker{
		sub: sub,
		obj: obj,
		db:  db,
		http: &http.Client{
			Timeout: ClientTimeOut,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: ClientTimeOut,
				}).DialContext,
			},
		},
	}
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.yaml", "Path to the configuration file")

	flag.Parse()

	conf, err := config.New(configPath)
	if err != nil {
		logger.Fatalf("config: %v", err)
	}

	logger.NewLogger(logger.LogLevel(conf.Log.Level), conf.Log.Mode == "dev")

	obj, err := minio.New(&conf.S3)
	if err != nil {
		logger.Fatalf("minio: %v", err)
	}

	sql, err := sqlstore.New(&conf.DB)
	if err != nil {
		logger.Fatalf("store: %v", err)
	}

	sub, err := nats.New(conf.EventAPI.Address(), conf.EventAPI.Enabled)
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
	w := newWorker(sub, obj, sql)
	if err := w.process(ctx); err != nil {
		logger.Fatalf("process: %v", err)
	}
}
