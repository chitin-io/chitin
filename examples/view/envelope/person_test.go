package person_test

import (
	"fmt"
	"log"

	"chitin.io/chitin/examples/view/envelope"
)

func Example() {
	input := []byte{
		// message kind: Person V2
		1,
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',
		// junk, e.g. the next message
		'x', 'y', 'z', 'z', 'y',
	}
	msg, err := person.NewMyProtocolView(input)
	if err != nil {
		log.Fatal(err)
	}
	switch view := msg.(type) {
	case *person.PersonV2View:
		fmt.Printf("%q is %d years old.\n", view.Name(), view.Age())
	default:
		log.Fatalf("unknown message type: %T", msg)
	}

	// Output:
	// "Jane" is 21 years old.
}
