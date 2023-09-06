package dummy

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/muyo/sno"
)

var ErrNoUser = errors.New("user does not exist")

type User struct {
	ID       string
	Name     string
	Email    string
	Gothic   []string
	Sessions []string
	Tokens   []string
}

type DB struct {
	users []User
}

func (d *DB) GetUserWithSession(_ context.Context, userID string, token string) (any, error) {

	for _, single := range d.users {
		if single.ID == userID {
			if !slices.Contains(single.Sessions, token) {
				return nil, fmt.Errorf("user found, but no sessions matching %s", token)
			}
			return single, nil
		}
	}

	return nil, ErrNoUser
}

func (d *DB) GetUser(_ context.Context, userID string) (any, error) {

	for _, single := range d.users {
		if single.ID == userID {
			return single, nil
		}
	}

	return nil, ErrNoUser
}
func (d *DB) GetUserIDFromToken(_ context.Context, token string) (string, error) {
	for _, single := range d.users {
		if slices.Contains(single.Tokens, token) {
			return single.ID, nil
		}
	}

	return "", ErrNoUser
}

func (d *DB) UpdateUserInfo(_ context.Context, userID string, email string, name string) error {

	for x, single := range d.users {
		if single.ID == userID {
			d.users[x].Email = email
			d.users[x].Name = name
			return nil
		}
	}

	return ErrNoUser
}

func (d *DB) CreateOrUpdateUser(_ context.Context, gothID, provider, email, name string) (string, error) {

	for x, single := range d.users {
		if slices.Contains(single.Gothic, gothID) {
			d.users[x].Email = email
			d.users[x].Name = name
			return single.ID, nil
		}
	}

	u := User{
		ID:       sno.New(0).String(),
		Name:     email,
		Email:    name,
		Gothic:   []string{gothID},
		Sessions: []string{},
		Tokens:   []string{sno.New(0).String()},
	}

	d.users = append(d.users, u)

	return u.ID, nil
}

func (d *DB) UserDisabled(_ context.Context, userID string) (bool, error) {

	for _, single := range d.users {
		if single.ID == userID {
			return false, nil
		}
	}

	return false, ErrNoUser
}

func (d *DB) GetUserID(_ context.Context, gothID string) (string, error) {

	for _, single := range d.users {
		if slices.Contains(single.Gothic, gothID) {
			return single.ID, nil
		}
	}

	return "", ErrNoUser
}

func (d *DB) AddSessionToUser(_ context.Context, gothID, session string) error {

	for x, single := range d.users {
		if slices.Contains(single.Gothic, gothID) {
			d.users[x].Sessions = append(single.Sessions, session)
			return nil
		}
	}

	return ErrNoUser
}
