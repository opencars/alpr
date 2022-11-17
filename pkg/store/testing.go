package store

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/opencars/alpr/pkg/domain/model"
)

func TestRecognition(t *testing.T) *model.Recognition {
	t.Helper()

	return &model.Recognition{
		ID:        uuid.NewV4().String(),
		ImageKey:  "plates/example.jpeg",
		Plate:     "AA1111AA",
		CreatedAt: time.Now().UTC(),
	}
}
