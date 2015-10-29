package schema

import (
	"bytes"
	"go/format"
	"io"
	"log"
)

type Schema interface {
	isSchema()

	Generate(dst io.Writer, pkg string) error
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

func (s *SchemaV1) Generate(dst io.Writer, pkg string) error {
	const debug = false

	var buf bytes.Buffer

	data := &genData{
		Package: pkg,
		Schema:  s,
	}
	if err := codegen.Execute(&buf, data); err != nil {
		return err
	}
	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		if debug {
			log.Printf("gofmt: %v", err)
			pretty = buf.Bytes()
		} else {
			return err
		}
	}
	if _, err := dst.Write(pretty); err != nil {
		return err
	}
	return nil
}

type Envelope map[uint]Versioned

type Message interface {
	isMessage()

	GetSlots() []Slot
	GetFields() []Field
}

type FixedMessage struct {
	WireFormat uint32
	Options    MessageOptions
	Slots      []Slot
}

func (*FixedMessage) isMessage() {}
func (*FixedMessage) isSlot()    {}
func (*FixedMessage) isField()   {}

func (m *FixedMessage) GetSlots() []Slot   { return m.Slots }
func (m *FixedMessage) GetFields() []Field { return nil }

type VarMessage struct {
	WireFormat uint32
	Options    MessageOptions
	Slots      []Slot
	Fields     []Field
}

func (*VarMessage) isMessage() {}
func (*VarMessage) isField()   {}

func (m *VarMessage) GetSlots() []Slot   { return m.Slots }
func (m *VarMessage) GetFields() []Field { return m.Fields }

type MessageOptions struct {
	Align uint32
}

type Slot struct {
	Name string
	Kind SlotType
}

type SlotType interface {
	isSlot()

	SlotSize() uint64
	GoType() string
	GoSlotGetter() string
}

type Field struct {
	Name    string
	Kind    FieldType
	Options FieldOptions
}

type FieldType interface {
	isField()

	GoType() string

	// Return the data type sufficient to hold data for fast field access.
	GoStateType() string

	// Return a template code snippet that assigns to the field given
	// in stateVar whatever is needed for later fast field access.
	// Local variable `data` is a `[]byte` with the unprocessed part
	// of the message, and is to be updated to point past the current
	// message.
	//
	// If field data is malformed, do `return nil, err`.
	GoFieldPrep(stateVar string) string

	GoFieldGetter(stateVar string) string
}

type FieldOptions struct {
	Align uint32
}
