package app_access_token

import (
	"time"
)

func Get(store Store, manager Manager, now time.Time) (string, error) {
	ok := false
	token, expires, err := store.Load()

	if err != nil {
		return "", err
	}

	if token == "" || now.After(expires) {
		goto request
	}

	ok, err = manager.Validate(token)

	if err != nil || ok {
		return token, err
	}

request:
	token, expires, err = manager.Request()

	if err != nil {
		return "", err
	}

	return token, store.Store(token, expires)
}
