package recognizer

import (
	"io"
)

// Recognizer ...
type Recognizer interface {
	Recognize(r io.Reader) ([]Result, error)
}
