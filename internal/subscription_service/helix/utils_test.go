package helix

import (
	"testing"

	"github.com/nicklaw5/helix/v2"
)

func TestErrorFromResponseNil(t *testing.T) {
	err := errorFromResponse(nil)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}
}

func TestErrorFromResponseOK(t *testing.T) {
	rc := &helix.ResponseCommon{
		Error:        "TestError",
		ErrorStatus:  456,
		ErrorMessage: "test error",
	}
	expected := "TestError (456): test error"
	err := errorFromResponse(rc)

	if err == nil {
		t.Errorf("err: nil; expected: %v", expected)
		return
	}

	msg := err.Error()

	if msg != expected {
		t.Errorf("msg: %v; expected: %v", err, expected)
	}
}
