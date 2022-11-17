package domain

import (
	"bytes"
	"context"
	"io"

	"github.com/opencars/alpr/pkg/domain/command"
	"github.com/opencars/alpr/pkg/domain/model"
	"github.com/opencars/alpr/pkg/domain/query"
)

type Store interface {
	Recognition() RecognitionRepository
}

type RecognitionRepository interface {
	Create(ctx context.Context, recognition *model.Recognition) error
	FindByPlate(ctx context.Context, plate string) ([]model.Recognition, error)
}

type InternalService interface {
	ListRecognitions(context.Context, *query.List) ([]model.Recognition, error)
}

type CustomerService interface {
	Recognize(context.Context, *command.Recognize) ([]model.Result, error)
}

// Recognizer is responsible for recognizing car plates from image as io.Reader.
type Recognizer interface {
	Recognize(r io.Reader) ([]model.Result, error)
}

// ObjectStore is responsible for uploading objects.
type ObjectStore interface {
	Put(ctx context.Context, key string, r *bytes.Reader) error
}

type Publisher interface {
	Publish(event *model.Event) error
}

type Subscriber interface {
	Subscribe(ctx context.Context) (<-chan *model.Event, error)
}

type Queue interface {
	Publisher
	Subscriber
}
