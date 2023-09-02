package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {

	check := require.New(t)

	for _, single := range List() {
		p, ok := FromSlug(single.String())
		check.True(ok)
		check.Equal(single, p)
		check.True(Validate(p))
		check.NotEmpty(p.Pretty())
		check.NotEmpty(p.String())
	}

	check.Equal(len(Slugs()), len(Public()))

	check.False(Validate(Provider{"abc", "abc"}))
	f, ok := FromSlug("abc")
	check.False(ok)
	check.False(Validate(f))

}
