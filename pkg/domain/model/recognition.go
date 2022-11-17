package model

import (
	"time"
)

const (
	// MaxImageSize equals to 5 MB.
	MaxImageSize = 5 << 20

	// ClientTimeOut equals to 10 seconds.
	ClientTimeOut = 10 * time.Second
)

type Recognition struct {
	ID        string    `db:"id"`
	ImageKey  string    `db:"image_key"`
	Plate     string    `db:"plate"`
	CreatedAt time.Time `db:"created_at"`
}
