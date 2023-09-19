package sqlite

import (
	"database/sql"
	"database/sql/driver"

	"modernc.org/sqlite"
)

const MemoryDB = "file:ent?mode=memory&cache=shared&_fk=1"

type sqlite3Driver struct {
	*sqlite.Driver
}

type sqlite3DriverConn interface {
	Exec(string, []driver.Value) (driver.Result, error)
}

func (d sqlite3Driver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	if _, err = conn.(sqlite3DriverConn).Exec("PRAGMA foreign_keys = ON;", nil); err != nil {
		_ = conn.Close()
	}
	return conn, nil
}

func init() {
	sql.Register("sqlite3", sqlite3Driver{Driver: &sqlite.Driver{}})
}
