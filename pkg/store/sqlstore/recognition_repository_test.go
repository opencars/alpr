package sqlstore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opencars/alpr/pkg/store"
	"github.com/opencars/alpr/pkg/store/sqlstore"
)

func TestRecognitionRepository_Create(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("recognitions")

	record := store.TestRecognition(t)
	assert.NoError(t, s.Recognition().Create(context.Background(), record))
	assert.NotEmpty(t, record.ID)
}

func TestRecognitionRepository_FindByPlate(t *testing.T) {
	s, teardown := sqlstore.TestDB(t, conf)
	defer teardown("recognitions")

	record := store.TestRecognition(t)
	assert.NoError(t, s.Recognition().Create(context.Background(), record))

	recognitions, err := s.Recognition().FindByPlate(context.Background(), record.Plate)
	assert.NoError(t, err)
	require.Len(t, recognitions, 1)
	assert.Equal(t, record.Plate, recognitions[0].Plate)
}
