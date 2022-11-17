package sqlstore

import (
	"context"

	"github.com/opencars/alpr/pkg/domain/model"
)

// RecognitionRepository is responsible for recognitions manipulation.
type RecognitionRepository struct {
	store *Store
}

// Create adds new record to the database.
func (r *RecognitionRepository) Create(ctx context.Context, recognition *model.Recognition) error {
	var id string
	err := r.store.db.GetContext(ctx, &id,
		`INSERT INTO recognitions (
			image_key, plate
		) VALUES (
			$1, $2
		) RETURNING id`,
		recognition.ImageKey, recognition.Plate,
	)

	if err != nil {
		return err
	}

	recognition.ID = id

	return nil
}

// FindByPlate returns list of records with specified plate.
func (r *RecognitionRepository) FindByPlate(ctx context.Context, plate string) ([]model.Recognition, error) {
	recognitions := make([]model.Recognition, 0)

	err := r.store.db.SelectContext(ctx, &recognitions,
		`SELECT id, image_key, plate, created_at
		FROM recognitions
		WHERE plate = $1`,
		plate,
	)

	if err != nil {
		return nil, err
	}

	return recognitions, nil
}
