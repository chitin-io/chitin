package person_test

//go:generate chitin -pkg person -i schema.chi -o chitin.gen.go

import (
	"fmt"

	"chitin.io/chitin/examples/maker/fields"
)

func Example() {
	m := person.NewPersonV2Maker()
	m.SetAge(15)
	m.SetSiblings(3)
	m.SetName("Jane")
	buf := m.Bytes()
	fmt.Printf("%q\n", buf)

	// Output:
	// "\x00\x0f\x00\x03\x05Jane"
}
