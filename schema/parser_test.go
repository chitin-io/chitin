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

func TestBadStart(t *testing.T) {
	_, err := schema.Parse(strings.NewReader(`foo`))
	if err == nil {
		t.Fatalf("expected an error")
	}
	if g, e := err.Error(), `1:1: unknown error: "foo"`; g != e {
		t.Fatalf("wrong error: %q != %q", g, e)
	}
}

func TestBadStart2(t *testing.T) {
	_, err := schema.Parse(strings.NewReader(`chitin foo`))
	if err == nil {
		t.Fatalf("expected an error")
	}
	if g, e := err.Error(), `1:8: unknown error: "foo"`; g != e {
		t.Fatalf("wrong error: %q != %q", g, e)
	}
}

func TestBadStart3(t *testing.T) {
	_, err := schema.Parse(strings.NewReader(`
foo`))
	if err == nil {
		t.Fatalf("expected an error")
	}
	if g, e := err.Error(), `2:1: unknown error: "foo"`; g != e {
		t.Fatalf("wrong error: %q != %q", g, e)
	}
}

func TestBadHeaderVersionBig(t *testing.T) {
	_, err := schema.Parse(strings.NewReader(`chitin v18446744073709551616
`))
	if err == nil {
		t.Fatalf("expected an error")
	}
	if g, e := err.Error(), `1:28: invalid uint32: value out of range: 18446744073709551616`; g != e {
		t.Fatalf("wrong error: %q != %q", g, e)
	}
}

func TestBadHeaderEOF(t *testing.T) {
	_, err := schema.Parse(strings.NewReader(`chitin v30`))
	if err == nil {
		t.Fatalf("expected an error")
	}
	if g, e := err.Error(), `1:8: unsupported schema version`; g != e {
		t.Fatalf("wrong error: %q != %q", g, e)
	}
}

func TestHeaderOnly(t *testing.T) {
	_, err := schema.Parse(strings.NewReader(`chitin v1
`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessage(t *testing.T) {
	_, err := schema.Parse(strings.NewReader(`
chitin v1

message Foo v13 {
	wire format: v1
}

message Bar v42 {
	wire format: v1
}
`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
