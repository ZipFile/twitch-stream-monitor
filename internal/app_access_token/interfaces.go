package app_access_token

import (
	"time"
)

type Store interface {
	Store(token string, expires time.Time) error
	Load() (string, time.Time, error)
}

type Manager interface {
	Validate(string) (bool, error)
	Request() (string, time.Time, error)
}
