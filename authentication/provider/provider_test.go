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
		check.NoError(p.Use("kbyuFDidLLm280LIwVFiazOqjO3ty8KH", "60Op4HFM0I8ajz0WdiStAbziZ-VFQttXuxixHHs2R7r7-CW8GR79l-mmLqMhc-Sa", "https://openidconnect.net/callback", "https://samples.auth0.com/.well-known/openid-configuration"))
	}

	check.Equal(len(Slugs()), len(Public()))

	check.False(Validate(Provider{"abc", "abc"}))
	f, ok := FromSlug("abc")
	check.False(ok)
	check.False(Validate(f))

	check.Error(OpenID.Use("123", "1234", "https://example.com", "https://example.com", "openid"))

}
