// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"chitin.io/chitin/schema"
)

func run() error {
	s := &schema.SchemaV1{
		Messages: map[schema.Versioned]schema.Message{
			{"Person", 1}: &schema.FixedMessage{
				WireFormat: 1,
				Slots: []schema.Slot{
					{
						Name: "age",
						Kind: schema.Uint16{},
					},
					{
						Name: "siblings",
						Kind: schema.Uint16{},
					},
				},
			},
		},
	}

	tmp, err := ioutil.TempFile(".", ".tmp-chitin-")
	if err != nil {
		log.Fatal(err)
	}
	doClose := true
	doRemove := true
	defer func() {
		if doClose {
			_ = tmp.Close()
		}
		if doRemove {
			_ = os.Remove(tmp.Name())
		}
	}()
	if err := s.Generate(tmp, "person"); err != nil {
		return fmt.Errorf("cannot generate wire format library: %v", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("error writing generated code: %v", err)
	}
	doClose = false
	if err := os.Rename(tmp.Name(), "chitin.gen.go"); err != nil {
		return fmt.Errorf("cannot rename temporary file: %v", err)
	}
	doRemove = false
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
