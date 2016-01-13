// This file is MANUALLY GENERATED and will be replaced by code
// generation once it starts working.

package person

import "encoding/binary"

func NewPersonV1Maker() *PersonV1Maker {
	// FixedMessage, size is well known
	maker := &PersonV1Maker{
		data: make([]byte, lenPersonV1Maker),
	}
	return maker
}

const (
	lenPersonV1Maker = 4
)

type PersonV1Maker struct {
	data []byte
}

func (m *PersonV1Maker) Bytes() []byte {
	return m.data
}

func (m *PersonV1Maker) SetAge(v uint16) {
	binary.BigEndian.PutUint16(m.data[0:2], v)
}

func (m *PersonV1Maker) SetSiblings(v uint16) {
	binary.BigEndian.PutUint16(m.data[2:4], v)
}
