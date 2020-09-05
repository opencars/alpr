package store

import (
	"time"
)

// Recognition ...
type Recognition struct {
	ID        string    `db:"id"`
	ImageKey  string    `db:"image_key"`
	Plate     string    `db:"plate"`
	CreatedAt time.Time `db:"created_at"`
}
