package person_test

//go:generate chitin -pkg person -i schema.chi -o chitin.gen.go

import (
	"fmt"

	"chitin.io/chitin/examples/maker/slots"
)

func Example() {
	m := person.NewPersonV1Maker()
	m.SetAge(15)
	m.SetSiblings(3)
	buf := m.Bytes()
	fmt.Printf("% 02x\n", buf)

	// Output:
	// 00 0f 00 03
}
