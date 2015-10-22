// This file is MANUALLY GENERATED and will be replaced by code
// generation once it starts working.

package person

import (
	"encoding/binary"
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

// Name returns a view of the name. Caller must not keep references to
// it past the lifetime of the view.
//
// If the message is truncated, returns an empty string.
func (v *PersonV2View) Name() string {
	l, n := varuint.Uint64(v.data[slotsLenPersonV2View:])
	if n < 0 {
		return ""
	}
	if l == 0 {
		panic("TODO padding not handled yet")
	}
	l--
	const maxInt = int(^uint(0) >> 1)
	if l > uint64(maxInt) {
		// technically, it has to be truncated because it wouldn't fit
		// in memory ;)
		return ""
	}
	end := slotsLenPersonV2View + n + int(l)
	if end > len(v.data) {
		return ""
	}
	var s string
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))
	p.Data = uintptr(unsafe.Pointer(&v.data[slotsLenPersonV2View+n]))
	p.Len = int(l)
	return s
}

func NewMyProtocolView(data []byte) (MyProtocolMessage, error) {
	if len(data) == 0 {
		return nil, chitin.ErrWrongSize
	}
	kind, n := varuint.Uint64(data)
	if n < 0 {
		return nil, chitin.ErrWrongSize
	}
	data = data[n:]
	switch kind {
	default:
		return nil, chitin.ErrUnknownMessageKind
	case 0:
		return nil, chitin.ErrIsPadding
	case 1:
		return NewPersonV2View(data)
	}
}

type MyProtocolMessage interface {
	isMyProtocolMessage()
}

func (*PersonV2View) isMyProtocolMessage() {}
