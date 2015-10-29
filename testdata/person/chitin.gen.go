// This file is automatically generated, DO NOT EDIT

package person

import (
	"encoding/binary"
	"io"
	"reflect"
	"unsafe"

	"chitin.io/chitin"
	"github.com/dchest/varuint"
)

// use all packages to avoid errors
var (
	_ = io.ErrUnexpectedEOF
	_ reflect.StringHeader
	_ unsafe.Pointer
	_ = varuint.Uint64
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
	slotsLenPersonV2View  = 4
	numFieldsPersonV2View = 1
	minLenPersonV2View    = slotsLenPersonV2View + 1*numFieldsPersonV2View
)

type PersonV2View struct {
	data []byte
}

func (v *PersonV2View) Age() uint16 {
	data := v.data[0:2]
	return binary.BigEndian.Uint16(data)
}

func (v *PersonV2View) Siblings() uint16 {
	data := v.data[2:4]
	return binary.BigEndian.Uint16(data)
}

func (v *PersonV2View) Fields() (*PersonV2ViewFields, error) {
	f := &PersonV2ViewFields{}
	off := slotsLenPersonV2View

	// TODO this only really implements length-prefixed fields

loop:
	l, n := varuint.Uint64(v.data[off:])
	if n < 0 {
		return nil, io.ErrUnexpectedEOF
	}
	if l == 0 {
		// padding
		off += n
		goto loop
	}
	l--

	const maxInt = int(^uint(0) >> 1)
	if l > uint64(maxInt) {
		// technically, it has to be truncated because it wouldn't fit
		// in memory ;)
		return nil, io.ErrUnexpectedEOF
	}
	li := int(l)

	// TODO prevent overflow here
	end := slotsLenPersonV2View + n + li
	if end > len(v.data) {
		return nil, io.ErrUnexpectedEOF
	}

	low := off + n
	high := low + li
	data := v.data[low:high]

	{
		p := (*reflect.StringHeader)(unsafe.Pointer(&f.fieldName))
		p.Data = uintptr(unsafe.Pointer(&data[0]))
		p.Len = len(data)
	}
	off = high

	return f, nil
}

type PersonV2ViewFields struct {
	fieldName string
}

func (f *PersonV2ViewFields) Name() string {
	return f.fieldName
}
