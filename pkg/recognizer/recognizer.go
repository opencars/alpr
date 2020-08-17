package recognizer

import (
	"io"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
)

type Recognizer interface {
	Recognize(r io.Reader) (*openalpr.AlprResults, error)
}
