// This file is MANUALLY GENERATED and will be replaced by code
// generation once it starts working.

package person

import (
	"encoding/binary"

	"github.com/dchest/varuint"
)

func NewPersonV2Maker() *PersonV2Maker {
	maker := &PersonV2Maker{}
	return maker
}

const (
	slotsLenPersonV2View  = 4
	numFieldsPersonV2View = 1
	minLenPersonV2View    = slotsLenPersonV2View + 1*numFieldsPersonV2View
)

type PersonV2Maker struct {
	slots     [slotsLenPersonV2View]byte
	fieldName string
}

func (m *PersonV2Maker) Bytes() []byte {
	// TODO what do we guarantee about immutability of return value?

	// TODO do this in just one allocation
	b := m.slots[:]

	var lb [varuint.MaxUint64Len]byte
	var ll int

	ll = varuint.PutUint64(lb[:], uint64(len(m.fieldName))+1)
	b = append(b, lb[:ll]...)
	b = append(b, m.fieldName...)

	return b
}

func (m *PersonV2Maker) SetAge(v uint16) {
	binary.BigEndian.PutUint16(m.slots[0:2], v)
}

func (m *PersonV2Maker) SetSiblings(v uint16) {
	binary.BigEndian.PutUint16(m.slots[2:4], v)
}

func (m *PersonV2Maker) SetName(v string) {
	m.fieldName = v
}
