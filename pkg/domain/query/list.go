package query

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/opencars/seedwork"
)

type List struct {
	Number string
}

func (q *List) Prepare() {}

func (q *List) Validate() error {
	return validation.ValidateStruct(q,
		validation.Field(
			&q.Number,
			validation.Required.Error(seedwork.Required),
			validation.Length(6, 18).Error(seedwork.Invalid),
		),
	)
}
