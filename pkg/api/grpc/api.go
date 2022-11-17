package grpc

import (
	"context"
	"net"

	"github.com/opencars/grpc/pkg/alpr"
	"google.golang.org/grpc"

	"github.com/opencars/alpr/pkg/domain"
)

// API represents gRPC API server.
type API struct {
	addr string
	s    *grpc.Server
	svc  domain.InternalService
}

func New(addr string, svc domain.InternalService) *API {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			RequestLoggingInterceptor,
		),
	}

	return &API{
		addr: addr,
		svc:  svc,
		s:    grpc.NewServer(opts...),
	}
}

func (a *API) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	alpr.RegisterServiceServer(a.s, &alprHandler{api: a})

	errors := make(chan error)
	go func() {
		errors <- a.s.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		a.s.GracefulStop()
		return <-errors
	case err := <-errors:
		return err
	}
}
