package track

import (
	"errors"
)

var (
	ErrNotFound = errors.New("track not found")
	ErrInvalid  = errors.New("invalid track")
)
