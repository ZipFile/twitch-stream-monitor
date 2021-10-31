package app_access_token

import (
	"testing"
	"time"
)

type StoreTest struct {
	Store Store
}

func (st StoreTest) FullCycle(t *testing.T) {
	token := "test"
	expires := time.Date(2021, time.October, 31, 11, 51, 0, 0, time.UTC)
	err := st.Store.Store(token, expires)

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
		return
	}

	loadedToken, loadedExpires, err := st.Store.Load()

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
		return
	}

	if loadedToken != token {
		t.Errorf("loadedToken: %v; expected: %v", loadedToken, token)
	}

	if !loadedExpires.Equal(expires) {
		t.Errorf("loadedExpires: %v; expected: %v", loadedExpires, expires)
	}
}

func (st StoreTest) Fresh(t *testing.T) {
	var initialExpires time.Time
	token, expires, err := st.Store.Load()

	if err != nil {
		t.Errorf("err: %v; expected: nil", err)
		return
	}

	if token != "" {
		t.Errorf("token: %v; expected: \"\"", token)
	}

	if !expires.Equal(initialExpires) {
		t.Errorf("expires: %v; expected: %v", expires, initialExpires)
	}
}
