package renderer

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

type Option interface {
	apply(*config) error
}

type option func(*config) error

func (o option) apply(c *config) error {
	return o(c)
}

type config struct {
	name  string
	files []string
	funcs template.FuncMap
}

var (
	ErrEmptyFunction  = errors.New("function is empty")
	ErrMissingTitle   = errors.New("missing function name")
	ErrFunctionExists = errors.New("function with title already exists")
	ErrEmptyFS        = errors.New("provided fs.FS is nil")
	ErrMissingFiles   = errors.New("missing names of files to use")
)

// Func adds a function to use in the templates. Title must be unique.
func Func(title string, function any) Option {
	return option(func(c *config) error {

		if function == nil {
			return ErrEmptyFunction
		}

		if title == "" {
			return ErrMissingTitle
		}

		keys := maps.Keys(c.funcs)

		for x := range keys {
			keys[x] = strings.ToLower(keys[x])
		}

		if slices.Contains(keys, strings.ToLower(title)) {
			return ErrFunctionExists
		}

		c.funcs[title] = function

		return nil

	})
}

// AddFiles to be used with the renderer.
// This function removes duplicates.
func AddFiles(files ...string) Option {
	return option(func(c *config) error {

		c.files = append(c.files, files...)

		c.files = slices.Compact(c.files)

		return nil
	})
}

func Name(name string) Option {
	return option(func(c *config) error {
		c.name = name
		return nil
	})
}

func New(files fs.FS, opts ...Option) (echo.Renderer, error) {

	if files == nil {
		return nil, ErrEmptyFS
	}

	var (
		err error
		c   = &config{
			funcs: make(template.FuncMap),
		}
	)

	for _, opt := range opts {
		err = errors.Join(err, opt.apply(c))
	}

	if err != nil {
		return nil, err
	}

	if len(c.files) < 1 {
		return nil, ErrMissingFiles
	}

	tmp, err := template.New(c.name).Funcs(c.funcs).ParseFS(files, c.files...)
	if err != nil {
		return nil, err
	}

	return &Renderer{tmp}, nil
}

type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
