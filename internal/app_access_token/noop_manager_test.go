package app_access_token

import (
	"errors"
	"testing"
)

func TestNoopManagerValidteOK(t *testing.T) {
	m := &NoopManager{Token: "test"}
	ok, err := m.Validate("test")

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if !ok {
		t.Errorf("ok: false; expected: true")
	}
}

func TestNoopManagerValidteError(t *testing.T) {
	testError := errors.New("test")
	m := &NoopManager{ValidateError: testError}
	ok, err := m.Validate("test")

	if err != testError {
		t.Errorf("err: %v; expected: testError", testError)
	}

	if ok {
		t.Errorf("ok: true; expected: false")
	}
}

func TestNoopManagerRequest(t *testing.T) {
	testError := errors.New("test")
	m := &NoopManager{Token: "test", Expires: expires, RequestError: testError}
	token, requestedExpires, err := m.Request()

	if err != testError {
		t.Errorf("err: %v; expected: testError", testError)
	}

	if token != "test" {
		t.Errorf("token: %v; expected: \"test\"", token)
	}

	if !requestedExpires.Equal(expires) {
		t.Errorf("token: %v; expected: %v", requestedExpires, expires)
	}
}
