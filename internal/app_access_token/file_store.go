package app_access_token

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

type fileStore struct {
	Path string
}

type value struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewFileStore(path string) Store {
	return &fileStore{path}
}

func (s *fileStore) Store(token string, expires time.Time) error {
	data, err := json.Marshal(value{token, expires})

	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.Path, data, 0600)
}

func (s *fileStore) Load() (string, time.Time, error) {
	var expires time.Time
	data, err := ioutil.ReadFile(s.Path)

	if errors.Is(err, os.ErrNotExist) {
		return "", expires, nil
	}

	var v value

	err = json.Unmarshal(data, &v)

	if err != nil {
		return "", expires, err
	}

	return v.Token, v.Expires, nil
}
