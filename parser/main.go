// chitin-parser project main.go
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"chitin.io/chitin/schema"
	"github.com/kr/pretty"
)

var (
	inputFile   = flag.String("i", "", "File to parse")
	outputFile  = flag.String("o", "", "File to generate")
	packageFile = flag.String("pkg", "main", "Package to generate")
	debug       = flag.Bool("d", false, "Print the schema")
)

func main() {
	flag.Parse()
	if *inputFile == "" {
		log.Fatal("No input file specified")
	}
	if *outputFile == "" {
		*outputFile = *inputFile + ".gen.go"
	}

	p := MakeParser()
	data, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	s, err := p.ParseString(string(data))
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		pretty.Print(s)
	}

	output, _ := os.Create(*outputFile)
	s.(schema.Schema).Generate(output, *packageFile)
}
