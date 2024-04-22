package broken

import (
	"errors"
	"io"
)

var _ io.Reader = &Reader{}

// Reader is an implementation of io.Reader that always malfunctions.
// Use only for testing purposes.
type Reader struct{}

func (r Reader) Read(_ []byte) (int, error) {
	return 0, errors.New("mock failure")
}
