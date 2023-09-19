package sqlite

import (
	"testing"

	"github.com/hcarriz/reverb/generated/ent/enttest"

	"entgo.io/ent/dialect"
	"github.com/stretchr/testify/require"
)

func TestSqlite(t *testing.T) {

	check := require.New(t)

	cl := enttest.Open(t, dialect.SQLite, MemoryDB)
	check.NotNil(cl)

	check.NoError(cl.Close())

}
