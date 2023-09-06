package dummy

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func TestDummy(t *testing.T) {

	check := require.New(t)

	d := &DB{}

	gothID := "12345"
	provider := "faux"

	id, err := d.CreateOrUpdateUser(nil, gothID, provider, fake.EmailAddress(), fake.FullName())
	check.NoError(err)
	check.NotEmpty(id)

	id2, err := d.CreateOrUpdateUser(nil, gothID, "2", fake.EmailAddress(), fake.FullName())
	check.NoError(err)
	check.Equal(id, id2)

	id3, err := d.GetUserID(nil, gothID)
	check.NoError(err)
	check.Equal(id, id3)

	_, err = d.GetUserID(nil, provider)
	check.Error(err)

	session := "67890"

	check.NoError(d.AddSessionToUser(nil, gothID, session))

	usr1, err := d.GetUserWithSession(nil, id3, session)
	check.NoError(err)
	check.NotEmpty(usr1)

	u2, ok := usr1.(User)
	check.True(ok)
	check.Equal(id3, u2.ID)

	usr2, err := d.GetUser(nil, id3)
	check.NoError(err)
	check.Equal(usr1, usr2)

}
