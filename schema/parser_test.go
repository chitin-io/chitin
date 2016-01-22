package schema

import "testing"

func TestFuzz1(t *testing.T) {
	parser := makeParser()

	parser.ParseString("chitin\tv1message 0\tv0{slots{}")

}

func TestFuzz2(t *testing.T) {
	parser := makeParser()

	parser.ParseString("chitin\tv1envelope\t0{}")
}

func TestFuzz3(t *testing.T) {
	parser := makeParser()

	parser.ParseString("chitin\nv1message\n0\nv0{slots{0\nstring")

}
