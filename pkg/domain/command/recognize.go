package command

import (
	"github.com/opencars/seedwork"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Recognize struct {
	URL string
}

func (c *Recognize) Prepare() {}

func (c *Recognize) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(
			&c.URL,
			validation.Required.Error(seedwork.Required),
			is.URL.Error(seedwork.Invalid),
		),
	)
}

func (c *Recognize) Event() {}
