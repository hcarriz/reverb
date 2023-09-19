package reverb

import (
	"testing"

	"entgo.io/ent/dialect"
	"github.com/hcarriz/reverb/generated/ent"
	"github.com/hcarriz/reverb/generated/ent/enttest"
	"github.com/hcarriz/reverb/sqlite"
	"github.com/stretchr/testify/require"

	_ "github.com/hcarriz/reverb/sqlite"
)

func TestReverb(t *testing.T) {

	check := require.New(t)

	e, err := New()
	check.NoError(err)
	check.NotEmpty(e)

}

func TestContextEnt(t *testing.T) {

	check := require.New(t)

	cl := enttest.Open(t, dialect.SQLite, sqlite.MemoryDB)

	check.NotNil(cl)

	AddEntToContext(ent.NewContext, cl)

}
