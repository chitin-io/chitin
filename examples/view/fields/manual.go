// This file is MANUALLY GENERATED and will be replaced by code
// generation once it starts working.

package person

import (
	"encoding/binary"
	"io"
	"reflect"
	"unsafe"

	"github.com/dchest/varuint"

	"chitin.io/chitin"
)

func NewPersonV2View(data []byte) (*PersonV2View, error) {
	if len(data) < minLenPersonV2View {
		return nil, chitin.ErrWrongSize
	}
	view := &PersonV2View{
		data: data,
	}
	return view, nil
}

const (
	maxInt = int(^uint(0) >> 1)

	slotsLenPersonV2View  = 4
	numFieldsPersonV2View = 1
	minLenPersonV2View    = slotsLenPersonV2View + 1*numFieldsPersonV2View
)

type PersonV2View struct {
	data []byte
}

func (v *PersonV2View) Age() uint16 {
	return binary.BigEndian.Uint16(v.data[0:2])
}

func (v *PersonV2View) Siblings() uint16 {
	return binary.BigEndian.Uint16(v.data[2:4])
}

// If the message is truncated, returns io.ErrUnexpectedEOF
func (v *PersonV2View) Fields() (*PersonV2ViewFields, error) {
	f := &PersonV2ViewFields{}

	l, n := varuint.Uint64(v.data[slotsLenPersonV2View:])
	if n < 0 {
		return nil, io.ErrUnexpectedEOF
	}
	if l == 0 {
		panic("TODO padding not handled yet")
	}
	l--
	if l > uint64(maxInt) {
		// technically, it has to be truncated because it wouldn't fit
		// in memory ;)
		return nil, io.ErrUnexpectedEOF
	}
	end := slotsLenPersonV2View + n + int(l)
	if end > len(v.data) {
		return nil, io.ErrUnexpectedEOF
	}
	p := (*reflect.StringHeader)(unsafe.Pointer(&f.name))
	p.Data = uintptr(unsafe.Pointer(&v.data[slotsLenPersonV2View+n]))
	p.Len = int(l)

	return f, nil
}

type PersonV2ViewFields struct {
	name string
}

// Name returns a view of the name. Caller must not keep references to
// it past the lifetime of the view.
func (f *PersonV2ViewFields) Name() string {
	return f.name
}
