package app_access_token

import (
	"time"
)

type NoopManager struct {
	Token         string
	Expires       time.Time
	ValidateError error
	RequestError  error
}

var expires time.Time = time.Date(3000, time.January, 1, 0, 0, 0, 0, time.UTC)

func (m *NoopManager) Validate(token string) (bool, error) {
	return m.Token == token, m.ValidateError
}

func (m *NoopManager) Request() (string, time.Time, error) {
	return m.Token, m.Expires, m.RequestError
}
