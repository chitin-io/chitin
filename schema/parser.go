package schema

import (
	"fmt"
	"io"
	"strconv"

	p "github.com/andyleap/parser"
)

func Parse(r io.ReadSeeker) (Schema, error) {
	p := makeParser()
	tmp, err := p.Parse(r)
	if err != nil {
		return nil, err
	}
	return tmp.(Schema), nil
}

func makeParser() *p.Grammar {
	ws := p.Mult(0, 0, p.Set("\\s"))
	rws := p.Mult(1, 0, p.Set("\\s"))

	schemaVersion := p.Or(p.Lit("v1"))
	intro := p.And(ws, p.Lit("chitin"), p.Require(rws, p.Tag("schemaVersion", schemaVersion)))

	versionIdent := p.And(p.Lit("v"), p.Tag("n", p.Mult(1, 0, p.Set("0-9"))))
	versionIdent.Node(func(m p.Match) (p.Match, error) {
		v, err := strconv.ParseUint(p.String(p.GetTag(m, "n")), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(v), nil
	})

	rawTypes := p.Or(
		p.Lit("uint8"),
		p.Lit("uint16"),
		p.Lit("uint32"),
		p.Lit("uint64"),
		p.Lit("int8"),
		p.Lit("int16"),
		p.Lit("int32"),
		p.Lit("int64"),
		p.Lit("float32"),
		p.Lit("float64"),
	)
	rawTypes.Node(func(m p.Match) (p.Match, error) {
		switch p.String(m) {
		case "uint16":
			return Uint16{}, nil
		}
		fmt.Printf("Bad Type: %v\n", m)
		return nil, fmt.Errorf("Type %s not supported", m)
	})

	aliasTypes := p.Or(
		p.Lit("byte"),
		p.Lit("string"),
	)
	aliasTypes.Node(func(m p.Match) (p.Match, error) {
		switch m.(p.MatchString) {
		case "string":
			return String{}, nil
		}
		fmt.Printf("Bad Type: %v\n", m)
		return nil, fmt.Errorf("Type %s not supported", m)
	})

	types := &p.Grammar{}

	fixedArrayType := p.And(
		p.Lit("["), p.Tag("n", p.Mult(1, 0, p.Set("\\d"))), p.Lit("]"),
		p.Tag("type", types),
	)
	fixedArrayType.Node(func(m p.Match) (p.Match, error) {
		return nil, fmt.Errorf("Type %s not supported", p.String(m))
	})

	varArrayType := p.And(
		p.Lit("[]"),
		p.Tag("type", types),
	)
	varArrayType.Node(func(m p.Match) (p.Match, error) {
		return nil, fmt.Errorf("Type %s not supported", p.String(m))
	})

	types.Set(p.Or(varArrayType, fixedArrayType, aliasTypes, rawTypes))

	wireFormat := p.And(p.Lit("wire format:"), p.Require(ws, p.Tag("wireFormat", versionIdent)))

	slotItemFormat := p.And(
		p.Tag("name", p.Mult(1, 0, p.Set("\\w"))), p.Require(ws, p.Tag("type", types), ws))
	slotItemFormat.Node(func(m p.Match) (p.Match, error) {
		slottype, ok := p.GetTag(m, "type").(SlotType)
		if !ok {
			return nil, fmt.Errorf("%T is not a SlotType", p.GetTag(m, "type"))
		}
		slot := Slot{
			Name: p.String(p.GetTag(m, "name")),
			Kind: slottype,
		}
		return p.TagMatch("slot", slot), nil
	})

	slotFormat := p.And(
		p.Lit("slots"), p.Require(ws,
			p.Lit("{"), ws,
			p.Mult(0, 0, p.And(slotItemFormat)),
			p.Lit("}"),
		))

	fieldItemFormat := p.And(
		p.Tag("name", p.Mult(1, 0, p.Set("\\w"))), p.Require(ws, p.Tag("type", types), ws))
	fieldItemFormat.Node(func(m p.Match) (p.Match, error) {
		field := Field{
			Name: p.String(p.GetTag(m, "name")),
			Kind: p.GetTag(m, "type").(FieldType),
		}
		return p.TagMatch("field", field), nil
	})

	fieldFormat := p.And(
		p.Lit("fields"), p.Require(ws,
			p.Lit("{"), ws,
			p.Mult(0, 0, p.And(fieldItemFormat)),
			p.Lit("}"),
		))

	messageParts := p.Or(wireFormat, slotFormat, fieldFormat)

	versionedName := p.And(
		p.Tag("name", p.Mult(1, 0, p.Set("\\w"))), rws,
		p.Tag("msgVersion", versionIdent),
	)
	versionedName.Node(func(m p.Match) (p.Match, error) {
		name := p.String(p.GetTag(m, "name"))
		version := p.GetTag(m, "msgVersion").(uint32)
		return Versioned{
			Name:    name,
			Version: version,
		}, nil
	})

	messageBody := p.And(p.Mult(1, 0, p.And(messageParts, ws)), ws)
	messageBody.Node(func(m p.Match) (p.Match, error) {

		wireFormat, _ := p.GetTag(m, "wireFormat").(uint32)

		slots := []Slot{}
		for _, slot := range p.GetTags(m, "slot") {
			slots = append(slots, slot.(Slot))
		}
		fields := []Field{}
		for _, field := range p.GetTags(m, "field") {
			fields = append(fields, field.(Field))
		}
		if len(fields) == 0 {
			msg := &FixedMessage{
				WireFormat: wireFormat,
				Slots:      slots,
			}
			return msg, nil
		}
		msg := &VarMessage{
			WireFormat: wireFormat,
			Slots:      slots,
			Fields:     fields,
		}
		return msg, nil
	})

	message := p.And(
		p.Lit("message"), p.Require(rws,
			p.Tag("name", versionedName), ws,
			p.Lit("{"), ws,
			p.Tag("body", messageBody), ws,
			p.Lit("}"),
		))

	envelopeMessages := p.And(
		p.Lit("messages"), p.Require(ws, p.Lit("{"), ws,
			p.Mult(1, 0, p.Tag("mapping", p.And(
				p.Tag("key", p.Mult(1, 0, p.Set("\\d"))), ws,
				p.Lit(":"), ws,
				p.Tag("value", versionedName), ws,
			))),
			p.Lit("}"),
		))
	envelopeMessages.Node(func(m p.Match) (p.Match, error) {
		messages := make(Envelope)
		for _, msg := range p.GetTags(m, "mapping") {
			key, _ := strconv.ParseUint(p.String(p.GetTag(msg, "key")), 10, 64)
			messages[uint(key)] = p.GetTag(msg, "value").(Versioned)
		}
		return messages, nil
	})

	envelopeParts := p.Or(p.Tag("body", envelopeMessages))

	envelope := p.And(
		p.Lit("envelope"), p.Require(rws, p.Tag("name", p.Mult(1, 0, p.Set("\\w"))), ws,
			p.Lit("{"), ws,
			p.Mult(0, 0, p.And(envelopeParts, ws)),
			p.Lit("}"),
		))

	schemaGrammar := p.And(
		intro, ws,
		p.Mult(0, 0, p.Or(
			p.And(p.Tag("message", message), ws),
			p.And(p.Tag("envelope", envelope), ws),
		)),
	)
	schemaGrammar.Node(func(m p.Match) (p.Match, error) {

		messages := make(map[Versioned]Message)
		for _, msg := range p.GetTags(m, "message") {
			name := p.GetTag(msg, "name").(Versioned)
			messages[name] = p.GetTag(msg, "body").(Message)
		}

		envelopes := make(map[string]Envelope)
		for _, msg := range p.GetTags(m, "envelope") {
			name := p.String(p.GetTag(msg, "name"))
			var e Envelope
			ok := false
			if e, ok = p.GetTag(msg, "body").(Envelope); !ok {
				return nil, fmt.Errorf("Envelope %s has no mapping!", name)
			}
			envelopes[name] = e
		}

		s := &SchemaV1{
			Messages:  messages,
			Envelopes: envelopes,
		}

		return s, nil
	})

	return schemaGrammar
}
