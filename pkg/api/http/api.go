package http

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/objectstore"
	"github.com/opencars/alpr/pkg/recognizer"
)

// Start starts the server with specified store.
func Start(ctx context.Context, addr string, conf *config.Server, rec recognizer.Recognizer, objStore objectstore.ObjectStore) error {
	s := newServer(rec, objStore)

	srv := http.Server{
		Addr:           addr,
		Handler:        handlers.LoggingHandler(os.Stdout, handlers.ProxyHeaders(s)),
		ReadTimeout:    conf.ReadTimeout.Duration,
		WriteTimeout:   conf.WriteTimeout.Duration,
		IdleTimeout:    conf.IdleTimeout.Duration,
		MaxHeaderBytes: 1 << 20,
	}

	errs := make(chan error)
	go func() {
		errs <- srv.ListenAndServe()
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		ctxShutDown, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout.Duration)
		defer cancel()

		err := srv.Shutdown(ctxShutDown)
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	}
}
