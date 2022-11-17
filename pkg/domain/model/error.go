package model

import (
	"github.com/opencars/seedwork"
)

var (
	ErrUnknownContentType = seedwork.NewError("request.url.content_type_invalid")
	ErrImageTooLarge      = seedwork.NewError("request.url.image_too_large")
)
