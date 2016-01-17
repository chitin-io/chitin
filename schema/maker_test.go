package schema_test

import (
	"testing"

	"chitin.io/chitin/testdata/person"
)

func TestMakerBytes(t *testing.T) {
	maker := person.NewPersonV2Maker()
	maker.SetAge(21)
	maker.SetSiblings(3)
	maker.SetName("Jane")
	maker.SetPhone("555-42")
	got := maker.Bytes()
	want := []byte{
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
	if g, e := got, want; string(g) != string(e) {
		t.Errorf("wrong bytes: %q != %q", g, e)
	}
}
