package service

import (
	"context"

	"github.com/opencars/seedwork"

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
	if err := seedwork.ProcessQuery(q); err != nil {
		return nil, err
	}

	results, err := s.repo.FindByPlate(ctx, q.Number)
	if err != nil {
		return nil, err
	}

	return results, nil
}
