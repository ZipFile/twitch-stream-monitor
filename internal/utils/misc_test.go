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

func TestObfuscateUrl(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{
			value:    "http://example.com",
			expected: "http://example.com",
		},
		{
			value:    "https://43bb-1-2-3-4.ngrok.io/test",
			expected: "***",
		},
		{
			value:    "http://1.2.3.4/test",
			expected: "***",
		},
	}

	for _, test := range tests {
		value := ObfuscateUrl(test.value)

		if value != test.expected {
			t.Errorf("value: %v; expected: %v", value, test.expected)
		}
	}
}
