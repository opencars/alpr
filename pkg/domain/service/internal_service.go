package service

import (
	"context"

	"github.com/opencars/alpr/pkg/domain"
	"github.com/opencars/alpr/pkg/domain/model"
	"github.com/opencars/alpr/pkg/domain/query"
)

type InternalService struct {
	repo domain.RecognitionRepository
	obj  domain.ObjectStore
}

func NewInternalService(repo domain.RecognitionRepository, obj domain.ObjectStore) *InternalService {
	return &InternalService{
		repo: repo,
		obj:  obj,
	}
}

func (s *InternalService) ListRecognitions(ctx context.Context, q *query.List) ([]model.Recognition, error) {
	return nil, nil
}
