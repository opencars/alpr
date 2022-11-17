package service

import (
	"context"

	"github.com/opencars/alpr/pkg/domain"
	"github.com/opencars/alpr/pkg/domain/model"
)

type InternalService struct {
	recognizer domain.Recognizer
	obj        domain.ObjectStore
}

func NewInternalService(rec domain.Recognizer, obj domain.ObjectStore) *InternalService {
	return &InternalService{
		recognizer: rec,
		obj:        obj,
	}
}

func (svc *InternalService) ListRecognitions(ctx context.Context, number string) ([]model.Recognition, error) {

	return nil, nil
}
