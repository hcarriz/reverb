package viewer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViewer(t *testing.T) {

	check := require.New(t)

	ctx := context.Background()

	id := "test"

	check.Equal("", GetUserID[string](ctx))

	ctx = SetUserID(ctx, id)

	check.Equal(id, GetUserID[string](ctx))

	check.Equal("127.0.0.1", GetAddress(ctx))

	ip := "192.168.1.1"

	ctx = SetAddress(ctx, ip)

	check.Equal(ip, GetAddress(ctx))

}
