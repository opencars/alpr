package grpc

import (
	"context"

	"github.com/opencars/alpr/pkg/domain/query"
	"github.com/opencars/grpc/pkg/alpr"
)

type alprHandler struct {
	alpr.UnimplementedServiceServer
	api *API
}

func (h *alprHandler) FindByNumber(ctx context.Context, r *alpr.NumberRequest) (*alpr.Response, error) {
	q := query.List{
		Number: r.Number,
	}

	result, err := h.api.svc.ListRecognitions(ctx, &q)
	if err != nil {
		return nil, handleErr(err)
	}

	dtos := make([]*alpr.Recognition, 0, len(result))

	for i := range result {
		dtos = append(dtos, convert(&result[i]))
	}

	return &alpr.Response{
		Recognitions: dtos,
	}, nil
}
