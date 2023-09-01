package cors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	tests := []struct {
		name    string
		args    []Option
		wantErr bool
	}{
		{
			name:    "default",
			wantErr: false,
		},
		{
			name: "failing origins",
			args: []Option{
				Origins("https://google.com/testing"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			check := require.New(t)

			_, err := New(tt.args...)
			if tt.wantErr {
				check.Error(err)
			} else {
				check.NoError(err)
			}

		})
	}
}
