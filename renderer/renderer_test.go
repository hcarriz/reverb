package renderer

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

const filename = "index.html"

func files() fs.FS {
	return fstest.MapFS{
		filename: {
			Data: []byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hello, World!</title>
</head>
<body>
    <p>Hello, World!</p> 
</body>
</html>`),
		},
	}
}

func TestRenderer(t *testing.T) {

	var (
		err   error
		check = require.New(t)
		e     = echo.New()
	)

	f := files()

	e.Renderer, err = New(f, AddFiles(filename))
	check.NoError(err)
	check.NotNil(e.Renderer)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	check.NoError(ctx.Render(http.StatusOK, filename, nil))

	defer rec.Result().Body.Close()

	result, err := io.ReadAll(rec.Result().Body)
	check.NoError(err)
	check.NotEmpty(result)

	raw, err := f.Open(filename)
	check.NoError(err)

	defer raw.Close()

	b, err := io.ReadAll(raw)
	check.NoError(err)

	check.Equal(result, b)

}

func TestNew(t *testing.T) {
	type args struct {
		files fs.FS
		opts  []Option
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		specificErr error
	}{
		// TODO: Add test cases.
		{
			name:        "empty fs",
			wantErr:     true,
			specificErr: ErrEmptyFS,
		},
		{
			name: "fs, but no functons",
			args: args{
				files: files(),
				opts:  []Option{},
			},
			wantErr: true,
		},
		{
			name: "fs and filename",
			args: args{
				files: files(),
				opts: []Option{
					AddFiles(filename),
				},
			},
			wantErr: false,
		},
		{
			name: "fs, filename, and function",
			args: args{
				files: files(),
				opts: []Option{
					AddFiles(filename),
					Func("now", time.Now),
				},
			},
			wantErr: false,
		},
		{
			name: "fs, filename, and function without title",
			args: args{
				files: files(),
				opts: []Option{
					AddFiles(filename),
					Func("", time.Now),
				},
			},
			wantErr:     true,
			specificErr: ErrMissingTitle,
		},
		{
			name: "fs, filename, and function addition without function",
			args: args{
				files: files(),
				opts: []Option{
					AddFiles(filename),
					Func("empty", nil),
				},
			},
			wantErr:     true,
			specificErr: ErrEmptyFunction,
		},
		{
			name: "fs, filename, and duplicate functions",
			args: args{
				files: files(),
				opts: []Option{
					AddFiles(filename),
					Func("now", time.Now),
					Func("now", time.Now),
				},
			},
			wantErr:     true,
			specificErr: ErrFunctionExists,
		},
		{
			name: "fs, filename, and duplicate functions with different titles",
			args: args{
				files: files(),
				opts: []Option{
					AddFiles(filename),
					Func("Now", time.Now),
					Func("nOW", time.Now),
				},
			},
			wantErr:     true,
			specificErr: ErrFunctionExists,
		},
		{
			name: "adding file that doesn't exist",
			args: args{
				files: files(),
				opts: []Option{
					AddFiles(filename, "false"),
				},
			},
			wantErr: true,
		},
		{
			name: "all options",
			args: args{
				files: files(),
				opts: []Option{
					Name("name"),
					AddFiles(filename),
					Func("now", time.Now),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			check := require.New(t)

			_, err := New(tt.args.files, tt.args.opts...)

			switch tt.wantErr {
			case true:
				check.Error(err)
				if tt.specificErr != nil {
					check.EqualError(err, tt.specificErr.Error())
				}
			case false:
				check.NoError(err)
			}

		})
	}
}
