package app_access_token

import (
	"errors"
	"testing"
	"time"
)

func TestGetLoadError(t *testing.T) {
	testError := errors.New("test")
	store := &NoopStore{LoadError: testError}
	now := time.Date(2021, time.October, 31, 12, 0, 0, 0, time.UTC)
	_, err := Get(store, nil, now)

	if err != testError {
		t.Errorf("err: %v; expected: testError", err)
	}
}

func TestGetValidateError(t *testing.T) {
	testError := errors.New("test")
	store := &NoopStore{
		Token:   "test",
		Expires: time.Date(2021, time.October, 31, 14, 0, 0, 0, time.UTC),
	}
	now := time.Date(2021, time.October, 31, 12, 0, 0, 0, time.UTC)
	manager := &NoopManager{ValidateError: testError}
	_, err := Get(store, manager, now)

	if err != testError {
		t.Errorf("err: %v; expected: testError", err)
	}
}

func TestGetRequestError(t *testing.T) {
	testError := errors.New("test")
	store := &NoopStore{
		Token:   "test",
		Expires: time.Date(2021, time.October, 31, 14, 0, 0, 0, time.UTC),
	}
	now := time.Date(2021, time.October, 31, 12, 0, 0, 0, time.UTC)
	manager := &NoopManager{RequestError: testError}
	_, err := Get(store, manager, now)

	if err != testError {
		t.Errorf("err: %v; expected: testError", err)
	}
}

func TestGetValid(t *testing.T) {
	store := &NoopStore{
		Token:   "test",
		Expires: time.Date(2021, time.October, 31, 14, 0, 0, 0, time.UTC),
	}
	now := time.Date(2021, time.October, 31, 12, 0, 0, 0, time.UTC)
	manager := &NoopManager{Token: "test"}
	token, err := Get(store, manager, now)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if token != "test" {
		t.Errorf("token: %v; expected: \"test\"", token)
	}
}

func TestGetExpired(t *testing.T) {
	store := &NoopStore{
		Token:   "test",
		Expires: time.Date(2021, time.October, 31, 10, 0, 0, 0, time.UTC),
	}
	now := time.Date(2021, time.October, 31, 12, 0, 0, 0, time.UTC)
	manager := &NoopManager{
		Token:   "test",
		Expires: time.Date(2021, time.October, 31, 14, 0, 0, 0, time.UTC),
	}
	token, err := Get(store, manager, now)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if token != "test" {
		t.Errorf("token: %v; expected: \"test\"", token)
	}

	_, expires, _ := store.Load()

	if expires != manager.Expires {
		t.Errorf("expires: %v; expected: %v", expires, manager.Expires)
	}
}

func TestGetMissing(t *testing.T) {
	store := &NoopStore{}
	now := time.Date(2021, time.October, 31, 12, 0, 0, 0, time.UTC)
	manager := &NoopManager{
		Token:   "test",
		Expires: time.Date(2021, time.October, 31, 14, 0, 0, 0, time.UTC),
	}
	token, err := Get(store, manager, now)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
	}

	if token != "test" {
		t.Errorf("token: %v; expected: \"test\"", token)
	}

	_, expires, _ := store.Load()

	if expires != manager.Expires {
		t.Errorf("expires: %v; expected: %v", expires, manager.Expires)
	}
}
