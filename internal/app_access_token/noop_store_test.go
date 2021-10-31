package app_access_token

import (
	"errors"
	"testing"
	"time"
)

func TestNoopStoreStoreError(t *testing.T) {
	testError := errors.New("test")
	s := &NoopStore{StoreError: testError}
	err := s.Store("", time.Time{})

	if err != testError {
		t.Errorf("err: %v; expected: testError", err)
		return
	}
}

func TestNoopStoreLoadError(t *testing.T) {
	testError := errors.New("test")
	s := &NoopStore{LoadError: testError}
	_, _, err := s.Load()

	if err != testError {
		t.Errorf("err: %v; expected: testError", err)
		return
	}
}

func TestNoopStoreFullCycle(t *testing.T) {
	s := &NoopStore{}
	(StoreTest{s}).FullCycle(t)
}

func TestNoopStoreFresh(t *testing.T) {
	s := &NoopStore{}
	(StoreTest{s}).Fresh(t)
}
