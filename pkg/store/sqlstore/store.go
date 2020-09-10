package sqlstore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/store"
)

// Store is postgres wrapper for store.Store.
type Store struct {
	db *sqlx.DB

	recognitionRepository *RecognitionRepository
}

// Recognition is responsible for recognitions manipulation.
func (s *Store) Recognition() store.RecognitionRepository {
	if s.recognitionRepository == nil {
		s.recognitionRepository = &RecognitionRepository{
			store: s,
		}
	}

	return s.recognitionRepository
}

// New returns new instance of store.
func New(conf *config.Database) (*Store, error) {
	info := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		conf.Host, conf.Port, conf.User, conf.Name, conf.SSLMode, conf.Password,
	)

	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return &Store{
		db: db,
	}, nil
}
