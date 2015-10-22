package schema

type Uint16 struct {
}

func (Uint16) isSlot()  {}
func (Uint16) isField() {}

type Byte struct {
}

func (Byte) isSlot()  {}
func (Byte) isField() {}

type Array struct {
	Length   uint64
	ItemKind SlotType
}

func (Array) isSlot()  {}
func (Array) isField() {}

type VarArray struct {
	ItemKind SlotType
}

func (VarArray) isField() {}

type String struct{}

func (String) isField() {}
