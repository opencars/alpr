package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/opencars/alpr/pkg/recognizer"

	"github.com/gorilla/handlers"
)

// Start starts the server with specified store.
func Start(ctx context.Context, addr string, recognizer recognizer.Recognizer) error {
	s := newServer(recognizer)
	srv := http.Server{
		Addr:    addr,
		Handler: handlers.LoggingHandler(os.Stdout, handlers.ProxyHeaders(s)),
	}

	errs := make(chan error)
	go func() {
		errs <- srv.ListenAndServe()
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer func() {
			cancel()
		}()

		err := srv.Shutdown(ctxShutDown)
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	}
}
