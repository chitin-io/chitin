// This file is MANUALLY GENERATED and will be replaced by code
// generation once it starts working.

package person

import (
	"encoding/binary"

	"chitin.io/chitin"
)

func NewPersonV1View(data []byte) (*PersonV1View, error) {
	// FixedMessage, size is well known
	if len(data) != lenPersonV1View {
		return nil, chitin.ErrWrongSize
	}
	view := &PersonV1View{
		data: data,
	}
	return view, nil
}

const (
	lenPersonV1View = 4
)

type PersonV1View struct {
	data []byte
}

func (v *PersonV1View) Age() uint16 {
	return binary.BigEndian.Uint16(v.data[0:2])
}

func (v *PersonV1View) Siblings() uint16 {
	return binary.BigEndian.Uint16(v.data[2:4])
}
