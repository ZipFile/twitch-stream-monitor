package app_access_token

import (
	"io/ioutil"
	"path"
	"testing"
	"time"
)

func TestFileStoreNew(t *testing.T) {
	s := NewFileStore("test").(*fileStore)

	if s.Path != "test" {
		t.Errorf("s.Path: %v; expected: \"test\"", s.Path)
	}
}

func TestFileStoreStore(t *testing.T) {
	f := path.Join(t.TempDir(), "test")
	s := NewFileStore(f)

	err := s.Store("test", time.Date(2021, time.October, 31, 11, 51, 0, 0, time.UTC))

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
		return
	}

	data, err := ioutil.ReadFile(f)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
		return
	}

	expected := `{"token":"test","expires":"2021-10-31T11:51:00Z"}`

	if string(data) != expected {
		t.Errorf("string(data): %v; expected: %v", string(data), expected)
	}
}

func TestFileStoreLoad(t *testing.T) {
	f := path.Join(t.TempDir(), "test")
	s := NewFileStore(f)
	data := `{"token":"test","expires":"2021-10-31T11:51:00Z"}`
	err := ioutil.WriteFile(f, []byte(data), 0644)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
		return
	}

	token, expires, err := s.Load()

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
		return
	}

	if token != "test" {
		t.Errorf("token: %v; expected: \"test\"", token)
	}

	expectedExpires := time.Date(2021, time.October, 31, 11, 51, 0, 0, time.UTC)

	if !expires.Equal(expectedExpires) {
		t.Errorf("expires: %v; expected: %v", expires, expectedExpires)
	}
}

func TestFileStoreFullCycle(t *testing.T) {
	s := NewFileStore(path.Join(t.TempDir(), "test"))
	(StoreTest{s}).FullCycle(t)
}

func TestFileStoreFresh(t *testing.T) {
	s := NewFileStore(path.Join(t.TempDir(), "test"))
	(StoreTest{s}).Fresh(t)
}
