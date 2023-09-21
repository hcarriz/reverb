package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {

	check := require.New(t)

	original := "the plain text password"

	// Check that the password is hashed.
	out, err := Create(original)
	check.NoError(err)
	check.NotEmpty(out)
	check.NotEqual(original, out.String())
	check.NotEmpty(out.String())

	// Check that we can verify the password.
	ok, err := Compare(original, *out)
	check.NoError(err)
	check.True(ok)

	// Check that the same password hashed again will have different results.
	out2, err := Create(original)
	check.NoError(err)
	check.NotEqual(out.String(), out2.String())

	// Check that wrong passwords won't be accepted.
	ok, err = Compare("not the right password", *out2)
	check.NoError(err)
	check.False(ok)

	// Verify CompareString works.
	ok, err = Check(original, out.String())
	check.NoError(err)
	check.True(ok)

}
