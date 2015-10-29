package person_test

//go:generate go run gen.go

import (
	"fmt"
	"log"

	"chitin.io/chitin/examples/view/slots"
)

func Example() {
	input := []byte{
		// age
		0, 21,
		// siblings
		0, 3,
	}
	view, err := person.NewPersonV1View(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("age:", view.Age())
	fmt.Println("siblings:", view.Siblings())

	// Output:
	// age: 21
	// siblings: 3
}
