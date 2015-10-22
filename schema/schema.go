package schema

type Schema interface {
	isSchema()
}

type Versioned struct {
	Name    string
	Version uint32
}

type SchemaV1 struct {
	Envelopes map[string]Envelope
	Messages  map[Versioned]Message
}

func (*SchemaV1) isSchema() {}

type Envelope map[uint]Versioned

type Message interface {
	isMessage()
}

type FixedMessage struct {
	WireFormat uint32
	Options    MessageOptions
	Slots      []Slot
}

func (*FixedMessage) isMessage() {}
func (*FixedMessage) isSlot()    {}
func (*FixedMessage) isField()   {}

type VarMessage struct {
	WireFormat uint32
	Options    MessageOptions
	Slots      []Slot
	Fields     []Field
}

func (*VarMessage) isMessage() {}
func (*VarMessage) isField()   {}

type MessageOptions struct {
	Align uint32
}

type Slot struct {
	Name string
	Kind SlotType
}

type SlotType interface {
	isSlot()
}

type Field struct {
	Name    string
	Kind    FieldType
	Options FieldOptions
}

type FieldType interface {
	isField()
}

type FieldOptions struct {
	Align uint32
}
