package model

import (
	"github.com/opencars/seedwork"
)

var (
	ErrRequiredImageURL   = seedwork.NewError("")
	ErrInvalidImageURL    = seedwork.NewError("")
	ErrUnknownContentType = seedwork.NewError("")
	ErrImageTooLarge      = seedwork.NewError("")
)
