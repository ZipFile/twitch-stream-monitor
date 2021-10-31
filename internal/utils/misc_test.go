package utils

import (
	"testing"
)

func TestMakeAddr(t *testing.T) {
	addr := MakeAddr("localhost", 9000)

	if addr != "localhost:9000" {
		t.Errorf("addr: %v; expected: \"localhost:9000\"", addr)
	}
}

func TestOrStrA(t *testing.T) {
	value := OrStr("a", "b")

	if value != "a" {
		t.Errorf("value: %v; expected: \"a\"", value)
	}
}

func TestOrStrB(t *testing.T) {
	value := OrStr("", "b")

	if value != "b" {
		t.Errorf("value: %v; expected: \"b\"", value)
	}
}
