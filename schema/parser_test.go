package schema_test

import (
	"strings"
	"testing"

	"chitin.io/chitin/schema"
)

func TestFuzz1(t *testing.T) {
	input := strings.NewReader("chitin\tv1message 0\tv0{slots{}")
	_, _ = schema.Parse(input)
}

func TestFuzz2(t *testing.T) {
	input := strings.NewReader("chitin\tv1envelope\t0{}")
	_, _ = schema.Parse(input)
}

func TestFuzz3(t *testing.T) {
	input := strings.NewReader("chitin\nv1message\n0\nv0{slots{0\nstring")
	_, _ = schema.Parse(input)
}
