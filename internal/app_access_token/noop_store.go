package app_access_token

import (
	"time"
)

type NoopStore struct {
	StoreError error
	LoadError  error
	Token      string
	Expires    time.Time
}

func (s *NoopStore) Store(token string, expires time.Time) error {
	s.Token = token
	s.Expires = expires

	return s.StoreError
}

func (s *NoopStore) Load() (string, time.Time, error) {
	return s.Token, s.Expires, s.LoadError
}
