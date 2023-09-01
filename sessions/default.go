package sessions

import (
	"context"
	"errors"
	"maps"
	"time"
)

func NewDB() *DB {
	return &DB{make(map[string]entry)}
}

type DB struct {
	entries map[string]entry
}

type entry struct {
	expires time.Time
	data    []byte
}

func (db *DB) Find(_ context.Context, token string) ([]byte, error) {

	result, ok := db.entries[token]
	if !ok {
		return nil, errors.New("invalid token")
	}

	return result.data, nil

}

func (db *DB) Delete(_ context.Context, token string) error {

	maps.DeleteFunc(db.entries, func(s string, e entry) bool {
		return s == token
	})

	return nil
}

func (db *DB) Add(_ context.Context, token string, data []byte, expires time.Time) error {

	db.entries[token] = entry{
		expires: expires,
		data:    data,
	}

	return nil
}

func (db *DB) All(_ context.Context) (map[string][]byte, error) {

	result := make(map[string][]byte, len(db.entries))

	for token, entry := range db.entries {
		result[token] = entry.data
	}

	return result, nil
}

func (db *DB) DeleteOld(_ context.Context, old time.Time) error {

	maps.DeleteFunc(db.entries, func(s string, e entry) bool {
		return e.expires.After(old)
	})

	return nil
}
