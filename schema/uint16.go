package schema

import "fmt"

type Uint16 struct {
}

var _ SlotType = Uint16{}
var _ FieldType = Uint16{}

func (Uint16) isSlot()  {}
func (Uint16) isField() {}

func (Uint16) SlotSize() uint64 {
	return 2
}

func (Uint16) GoType() string {
	return `uint16`
}

func (Uint16) GoSlotGetter() string {
	return `return binary.BigEndian.Uint16(data)`
}

func (Uint16) GoSlotSetter() string {
	return `binary.BigEndian.PutUint16(data, v)`
}

func (Uint16) GoStateType() string { return `uint16` }

func (Uint16) GoFieldPrep(stateVar string) string {
	return fmt.Sprintf(`
	v, n := varuint.Uint64(data)
	if n < 0 {
		return nil, io.ErrUnexpectedEOF
	}
	const maxUint16 = ^uint16(0)
	if v > maxUint16 {
		return nil, errors.New("value overflows uint16")
	}
	%s = uint16(v)
	data = data[n:]
`, stateVar)
}

func (Uint16) GoFieldGetter(stateVar string) string {
	return fmt.Sprintf(`return %s`, stateVar)
}

func (Uint16) GoFieldBytes(stateVar string) string {
	return fmt.Sprintf(`
	var lb [varuint.MaxUint64Len]byte
	var ll int

	ll = varuint.PutUint64(lb[:], %s)
	data = append(data, lb[:ll]...)
`, stateVar)
}
