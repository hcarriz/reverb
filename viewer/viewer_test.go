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

	fid, ok := GetUserID[string](ctx)
	check.False(ok)

	check.Equal("", fid)

	ctx = SetUserID(ctx, id)

	fid, ok = GetUserID[string](ctx)
	check.True(ok)

	check.Equal(id, fid)

	check.Equal("127.0.0.1", GetAddress(ctx))

	ip := "192.168.1.1"

	ctx = SetAddress(ctx, ip)

	check.Equal(ip, GetAddress(ctx))

	check.False(IsSystem(ctx))

	ctx = SetSystem(ctx)

	check.True(IsSystem(ctx))

}

func TestComplete(t *testing.T) {

	check := require.New(t)

	check.Equal(ContextIP.v, ContextIP.String())

}
