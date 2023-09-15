package faux

import (
	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

type Faux struct {
	name string
}

func (f *Faux) Name() string {
	return f.name
}

func (f *Faux) SetName(name string) {
	f.name = name
}

func (f *Faux) BeginAuth(state string) (goth.Session, error) {

	c := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL: "http://example.com/auth",
		},
	}
	url := c.AuthCodeURL(state)
	return &Session{
		ID:      "id",
		AuthURL: url,
	}, nil
}

func (f *Faux) UnmarshalSession(session string) (Session, error) {

	s := Session{}

	return s, nil
}

func (f *Faux) Debug(d bool) {}

func (f *Faux) RefreshToken(rt string) (*oauth2.Token, error) { return nil, nil }

func (f *Faux) RefreshTokenAvailable() bool { return false }

func New() *Faux {

	f := &Faux{
		name: "faux",
	}

	return f
}
