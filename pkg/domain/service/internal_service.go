package service

import (
	"context"
	"net/http"

	"github.com/opencars/alpr/pkg/domain/model"
)

type InternalService struct {
	client     *http.Client
	recognizer recognizer.Recognizer
	obj        objectstore.ObjectStore
	store      store.Store
	pub        queue.Publisher
}

func NewInternalService() *InternalService {
	return &InternalService{}
}

func (svc *InternalService) ListRecognitions(ctx context.Context, number string) ([]model.Recognition, error) {
	return nil, nil
}
