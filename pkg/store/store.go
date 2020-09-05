package store

import (
	"context"
)

type Store interface {
	Recognition() RecognitionRepository
}

type RecognitionRepository interface {
	Create(ctx context.Context, recognition *Recognition) error
	FindByPlate(ctx context.Context, plate string) ([]Recognition, error)
}
