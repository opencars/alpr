package recognizer

import (
	"errors"
)

var (
	// ErrFailedToLoad returned, if ALPR library failed to load.
	ErrFailedToLoad = errors.New("recognizer: failed to load")
)
