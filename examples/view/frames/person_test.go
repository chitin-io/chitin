package person_test

import (
	"fmt"
	"io"
	"log"

	"chitin.io/chitin"
	"chitin.io/chitin/examples/view/frames"
)

func Example_bytes() {
	input := []byte{
		// frame length +1
		10,
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',

		// frame length +1
		9,
		// age
		0, 30,
		// siblings
		0, 0,
		// len of name +1
		4,
		// name
		'J', 'o', 'e',
	}
	framed := chitin.NewFramedView(input)
	for {
		buf, err := framed.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		view, err := person.NewPersonV2View(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%q is %d years old.\n", view.Name(), view.Age())
	}

	// Output:
	// "Jane" is 21 years old.
	// "Joe" is 30 years old.
}

func Example_messages() {
	input := []byte{
		// frame length +1
		10,
		// age
		0, 21,
		// siblings
		0, 3,
		// len of name +1
		5,
		// name
		'J', 'a', 'n', 'e',

		// frame length +1
		9,
		// age
		0, 30,
		// siblings
		0, 0,
		// len of name +1
		4,
		// name
		'J', 'o', 'e',
	}
	framed := person.NewFramedPersonV2View(input)
	for {
		view, err := framed.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%q is %d years old.\n", view.Name(), view.Age())
	}

	// Output:
	// "Jane" is 21 years old.
	// "Joe" is 30 years old.
}
