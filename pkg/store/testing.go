package store

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestRecognition(t *testing.T) *Recognition {
	t.Helper()

	return &Recognition{
		ID:        uuid.NewV4().String(),
		ImageKey:  "plates/example.jpeg",
		Plate:     "AA1111AA",
		CreatedAt: time.Now().UTC(),
	}
}
