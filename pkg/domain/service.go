package domain

import (
	"context"

	"github.com/opencars/alpr/pkg/domain/model"
)

type InternalService interface {
	ListRecognitions(ctx context.Context, number string) ([]model.Recognition, error)
}
