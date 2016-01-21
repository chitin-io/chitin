// chitin-parser project main.go
package main

import (
	"fmt"
	"strconv"

	"chitin.io/chitin/schema"
	. "github.com/andyleap/parser"
)

func MakeParser() *Grammar {
	ws := Mult(0, 0, Set("\\s"))
	rws := Mult(1, 0, Set("\\s"))

	schemaVersion := Or(Lit("v1"))
	intro := And(ws, Lit("chitin"), Require(rws, Tag("schemaVersion", schemaVersion)))

	versionIdent := And(Lit("v"), Tag("n", Mult(1, 0, Set("0-9"))))
	versionIdent.Node(func(m Match) (Match, error) {
		v, err := strconv.ParseUint(String(GetTag(m, "n")), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(v), nil
	})

	rawTypes := Or(
		Lit("uint8"),
		Lit("uint16"),
		Lit("uint32"),
		Lit("uint64"),
		Lit("int8"),
		Lit("int16"),
		Lit("int32"),
		Lit("int64"),
		Lit("float32"),
		Lit("float64"),
	)
	rawTypes.Node(func(m Match) (Match, error) {
		switch String(m) {
		case "uint16":
			return schema.Uint16{}, nil
		}
		fmt.Printf("Bad Type: %v\n", m)
		return nil, fmt.Errorf("Type %s not supported", m)
	})

	aliasTypes := Or(
		Lit("byte"),
		Lit("string"),
	)
	aliasTypes.Node(func(m Match) (Match, error) {
		switch m.(MatchString) {
		case "string":
			return schema.String{}, nil
		}
		fmt.Printf("Bad Type: %v\n", m)
		return nil, fmt.Errorf("Type %s not supported", m)
	})

	types := &Grammar{}

	fixedArrayType := And(
		Lit("["), Tag("n", Mult(1, 0, Set("\\d"))), Lit("]"),
		Tag("type", types),
	)
	fixedArrayType.Node(func(m Match) (Match, error) {
		return nil, fmt.Errorf("Type %s not supported", String(m))
	})

	varArrayType := And(
		Lit("[]"),
		Tag("type", types),
	)
	varArrayType.Node(func(m Match) (Match, error) {
		return nil, fmt.Errorf("Type %s not supported", String(m))
	})

	types.Set(Or(varArrayType, fixedArrayType, aliasTypes, rawTypes))

	wireFormat := And(Lit("wire format:"), Require(ws, Tag("wireFormat", versionIdent)))

	slotItemFormat := And(
		Tag("name", Mult(1, 0, Set("\\w"))), Require(ws, Tag("type", types), ws))
	slotItemFormat.Node(func(m Match) (Match, error) {
		slot := schema.Slot{
			Name: String(GetTag(m, "name")),
			Kind: GetTag(m, "type").(schema.SlotType),
		}
		return slot, nil
	})

	slotFormat := And(
		Lit("slots"), Require(ws,
			Lit("{"), ws,
			Tag("slots", Mult(0, 0, And(slotItemFormat))),
			Lit("}"),
		))

	fieldItemFormat := And(
		Tag("name", Mult(1, 0, Set("\\w"))), Require(ws, Tag("type", types), ws))
	fieldItemFormat.Node(func(m Match) (Match, error) {
		field := schema.Field{
			Name: String(GetTag(m, "name")),
			Kind: GetTag(m, "type").(schema.FieldType),
		}
		return field, nil
	})

	fieldFormat := And(
		Lit("fields"), Require(ws,
			Lit("{"), ws,
			Tag("fields", Mult(0, 0, And(fieldItemFormat))),
			Lit("}"),
		))

	messageParts := Or(wireFormat, slotFormat, fieldFormat)

	versionedName := And(
		Tag("name", Mult(1, 0, Set("\\w"))), rws,
		Tag("msgVersion", versionIdent),
	)
	versionedName.Node(func(m Match) (Match, error) {
		name := String(GetTag(m, "name"))
		version := GetTag(m, "msgVersion").(uint32)
		return schema.Versioned{
			Name:    name,
			Version: version,
		}, nil
	})

	messageBody := And(Mult(1, 0, And(messageParts, ws)), ws)
	messageBody.Node(func(m Match) (Match, error) {

		wireFormat, _ := GetTag(m, "wireFormat").(uint32)

		slots := []schema.Slot{}
		for _, slot := range GetTag(m, "slots").(MatchTree) {
			slots = append(slots, slot.(MatchTree)[0].(schema.Slot))
		}
		fields := []schema.Field{}
		for _, field := range GetTag(m, "fields").(MatchTree) {
			fields = append(fields, field.(MatchTree)[0].(schema.Field))
		}
		if len(fields) == 0 {
			msg := &schema.FixedMessage{
				WireFormat: wireFormat,
				Slots:      slots,
			}
			return msg, nil
		}
		msg := &schema.VarMessage{
			WireFormat: wireFormat,
			Slots:      slots,
			Fields:     fields,
		}
		return msg, nil
	})

	message := And(
		Lit("message"), Require(rws,
			Tag("name", versionedName), ws,
			Lit("{"), ws,
			Tag("body", messageBody), ws,
			Lit("}"),
		))

	envelopeMessages := And(
		Lit("messages"), Require(ws, Lit("{"), ws,
			Mult(1, 0, Tag("mapping", And(
				Tag("key", Mult(1, 0, Set("\\d"))), ws,
				Lit(":"), ws,
				Tag("value", versionedName), ws,
			))),
			Lit("}"),
		))
	envelopeMessages.Node(func(m Match) (Match, error) {
		messages := make(schema.Envelope)
		for _, msg := range GetTags(m, "mapping") {
			key, _ := strconv.ParseUint(String(GetTag(msg, "key")), 10, 64)
			messages[uint(key)] = GetTag(msg, "value").(schema.Versioned)
		}
		return messages, nil
	})

	envelopeParts := Or(Tag("body", envelopeMessages))

	envelope := And(
		Lit("envelope"), Require(rws, Tag("name", Mult(1, 0, Set("\\w"))), ws,
			Lit("{"), ws,
			Mult(0, 0, And(envelopeParts, ws)),
			Lit("}"),
		))

	schemaGrammar := And(
		intro, ws,
		Mult(0, 0, Or(
			And(Tag("message", message), ws),
			And(Tag("envelope", envelope), ws),
		)),
	)
	schemaGrammar.Node(func(m Match) (Match, error) {

		messages := make(map[schema.Versioned]schema.Message)
		for _, msg := range GetTags(m, "message") {
			name := GetTag(msg, "name").(schema.Versioned)
			messages[name] = GetTag(msg, "body").(schema.Message)
		}

		envelopes := make(map[string]schema.Envelope)
		for _, msg := range GetTags(m, "envelope") {
			name := String(GetTag(msg, "name"))
			envelopes[name] = GetTag(msg, "body").(schema.Envelope)
		}

		s := &schema.SchemaV1{
			Messages:  messages,
			Envelopes: envelopes,
		}

		return s, nil
	})

	return schemaGrammar
}
