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

func (Uint16) GoStateType() string { return `uint16` }

func (Uint16) GoFieldPrep(stateVar string) string {
	// TODO not a length-prefixed field
	return fmt.Sprintf(`%s = binary.BigEndian.Uint16(data)`, stateVar)
}

func (Uint16) GoFieldGetter(stateVar string) string {
	return fmt.Sprintf(`return %s`, stateVar)
}
