package store

import (
	"context"

	"github.com/opencars/alpr/pkg/domain/model"
)

type Store interface {
	Recognition() RecognitionRepository
}

type RecognitionRepository interface {
	Create(ctx context.Context, recognition *model.Recognition) error
	FindByPlate(ctx context.Context, plate string) ([]model.Recognition, error)
}
