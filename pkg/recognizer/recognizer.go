package recognizer

import (
	"io"
)

// Recognizer is responsible for recognizing car plates from image as io.Reader.
type Recognizer interface {
	Recognize(r io.Reader) ([]Result, error)
}
