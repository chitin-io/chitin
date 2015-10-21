package chitin

import "errors"

var (
	// ErrWrongSize is the error returned when data is of an unexpected size.
	ErrWrongSize = errors.New("message length does not match")

	// ErrIsPadding is returned when an envelope contains no data.
	ErrIsPadding = errors.New("unexpected padding")

	// ErrUnknownMessageKind is the error returned when an envelope
	// contains an unrecognized message kind.
	ErrUnknownMessageKind = errors.New("unknown message kind")
)
