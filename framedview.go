package chitin

import (
	"io"

	"github.com/dchest/varuint"
)

// NewFramedView opens a view on the bytes given, providing access to
// framed contents as bytes.
func NewFramedView(data []byte) *FramedView {
	return &FramedView{data: data}
}

// FramedView is a view on framed bytes.
//
// See https://chitin.io/spec/v1/#f-e-m
type FramedView struct {
	data []byte
}

// Next returns the contents of the next frame.
//
// Returns io.EOF if there are no more frames. Returns
// io.ErrUnexpectedEOF on truncated data.
func (v *FramedView) Next() ([]byte, error) {
loop:
	if len(v.data) == 0 {
		return nil, io.EOF
	}
	l, n := varuint.Uint64(v.data)
	if n < 0 {
		return nil, io.ErrUnexpectedEOF
	}
	if l == 0 {
		goto loop
	}
	l--
	const maxInt = int(^uint(0) >> 1)
	if l > uint64(maxInt) {
		// technically, it has to be truncated because it wouldn't fit
		// in memory ;)
		return nil, io.ErrUnexpectedEOF
	}
	end := n + int(l)
	if end > len(v.data) {
		return nil, io.ErrUnexpectedEOF
	}
	b := v.data[n:end]
	v.data = v.data[end:]
	return b, nil
}
