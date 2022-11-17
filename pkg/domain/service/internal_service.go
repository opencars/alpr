package service

import (
	"context"

	"github.com/opencars/alpr/pkg/domain/model"
)

type InternalService struct {
	// client     *http.Client
	// recognizer domain.Recognizer
	// obj        domain.ObjectStore
	// store      domain.Store/
	// pub        domain.Publisher
}

func NewInternalService() *InternalService {
	return &InternalService{}
}

func (svc *InternalService) ListRecognitions(ctx context.Context, number string) ([]model.Recognition, error) {
	return nil, nil
}
