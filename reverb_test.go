package reverb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverb(t *testing.T) {

	check := require.New(t)

	e, err := New()
	check.NoError(err)
	check.NotEmpty(e)

}
