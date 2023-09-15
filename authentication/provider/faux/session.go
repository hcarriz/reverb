package faux

import (
	"encoding/json"

	"github.com/markbates/goth"
)

type Session struct {
	ID          string
	Name        string
	Email       string
	AuthURL     string
	AccessToken string
}

func (s *Session) GetAuthURL() (string, error) {
	return s.AuthURL, nil
}

func (s *Session) Marshal() string {

	if data, err := json.Marshal(s); err == nil {
		return string(data)
	}

	return ""
}

func (s *Session) Authorize(goth.Provider, goth.Params) (string, error) {
	return s.AccessToken, nil
}
