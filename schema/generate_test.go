package schema_test

import (
	"io"
	"testing"

	"chitin.io/chitin"
	"chitin.io/chitin/testdata/person"
)

func TestViewBadTruncatedSlots(t *testing.T) {
	// truncated so short not even all the slot data is present
	input := []byte{
		// age
		0, 21,
		// siblings
		0,
	}

	for i := len(input); i >= 0; i-- {
		view, err := person.NewPersonV2View(input)
		if g, e := err, chitin.ErrWrongSize; g != e {
			t.Errorf("wrong error for truncated input @%d: %q (%T) != %q (%T), view=%#v", i, g, g, e, e, view)
		}
	}
}

func TestViewSlot(t *testing.T) {
	input := []byte{
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',
		// len of phone +1
		7,
		// phone
		'5', '5', '5', '-', '4', '2',
	}
	view, err := person.NewPersonV2View(input)
	if err != nil {
		t.Fatalf("cannot open view: %v", err)
	}
	if g, e := view.Age(), uint16(21); g != e {
		t.Errorf("wrong age: %d != %d", g, e)
	}
	if g, e := view.Siblings(), uint16(3); g != e {
		t.Errorf("wrong siblings: %d != %d", g, e)
	}
}

func TestViewField(t *testing.T) {
	input := []byte{
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',
		// len of phone +1
		7,
		// phone
		'5', '5', '5', '-', '4', '2',
	}
	view, err := person.NewPersonV2View(input)
	if err != nil {
		t.Fatalf("cannot open view: %v", err)
	}
	fields, err := view.Fields()
	if err != nil {
		t.Fatalf("cannot access fields: %v", err)
	}
	if g, e := fields.Name(), `Jane`; g != e {
		t.Errorf("wrong name: %q != %q", g, e)
	}
	if g, e := fields.Phone(), `555-42`; g != e {
		t.Errorf("wrong name: %q != %q", g, e)
	}
}

func TestViewFieldPadding(t *testing.T) {
	input := []byte{
		// age
		0, 21,
		// siblings
		0, 3,
		// padding
		0, 0, 0,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',
		// padding
		0,
		// len of phone +1
		7,
		// phone
		'5', '5', '5', '-', '4', '2',
	}
	view, err := person.NewPersonV2View(input)
	if err != nil {
		t.Fatalf("cannot open view: %v", err)
	}
	fields, err := view.Fields()
	if err != nil {
		t.Fatalf("cannot access fields: %v", err)
	}
	if g, e := fields.Name(), `Jane`; g != e {
		t.Errorf("wrong name: %q != %q", g, e)
	}
}

func TestViewFieldBadTruncated(t *testing.T) {
	// truncated somewhere after slots
	input := []byte{
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',
		// len of phone +1
		7,
		// phone
		'5', '5', '5', '-', '4',
	}

	for i := len(input); i >= 4; i-- {
		view, err := person.NewPersonV2View(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		fields, err := view.Fields()
		if g, e := err, io.ErrUnexpectedEOF; g != e {
			t.Errorf("wrong error for truncated input @%d: %q (%T) != %q (%T), fields=%#v", i, g, g, e, e, fields)
		}
	}
}

func TestViewFieldBadTruncatedLen(t *testing.T) {
	input := []byte{
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',
		// len of phone +1, truncated varuint
		241,
	}
	view, err := person.NewPersonV2View(input)
	if err != nil {
		t.Fatalf("cannot open view: %v", err)
	}
	fields, err := view.Fields()
	if g, e := err, io.ErrUnexpectedEOF; g != e {
		t.Errorf("wrong error for large len: %q (%T) != %q (%T), fields=%#v", g, g, e, e, fields)
	}
}

func TestViewFieldBadLen(t *testing.T) {
	input := []byte{
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1, bad because it'll overflow (signed) int
		255,
		255, 255, 255, 255, 255, 255, 255, 255,
	}
	view, err := person.NewPersonV2View(input)
	if err != nil {
		t.Fatalf("cannot open view: %v", err)
	}
	fields, err := view.Fields()
	if g, e := err, io.ErrUnexpectedEOF; g != e {
		t.Errorf("wrong error for large len: %q (%T) != %q (%T), fields=%#v", g, g, e, e, fields)
	}
}
