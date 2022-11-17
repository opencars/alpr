package domain

import (
	"context"
	"io"

	"github.com/opencars/alpr/pkg/domain/model"
)

type InternalService interface {
	ListRecognitions(ctx context.Context, number string) ([]model.Recognition, error)
}

type CustomerService interface {
	Recognize(ctx context.Context, img string) ([]model.Recognition, error)
}

// Recognizer is responsible for recognizing car plates from image as io.Reader.
type Recognizer interface {
	Recognize(r io.Reader) ([]model.Result, error)
}
