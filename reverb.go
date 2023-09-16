package reverb

import (
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"slices"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type Auth interface {
	Authenticate(echo.Context) (context.Context, bool)
}

type pathfinder struct {
	method    string
	path      string
	handler   echo.HandlerFunc
	middlware []echo.MiddlewareFunc
}

type config struct {
	echo       *echo.Echo
	logger     *slog.Logger
	showBanner bool
	secrets    []string
	debug      bool
	session    *scs.SessionManager
}

// Errors
var (
	ErrEmptyPath         = errors.New("empty path")
	ErrEmptyRenderer     = errors.New("renderer is empty")
	ErrInvalidDuration   = errors.New("invalid duration")
	ErrInvalidHTTPMethod = errors.New("invalid http method")
	ErrInvalidLogger     = errors.New("invalid logger")
	ErrMissingRoutes     = errors.New("missing route")
)

// Option is a functional optional pattern
type Option interface {
	apply(*config) error
}

type Options []Option

func (o Options) apply(a *config) error {

	var err error

	for _, single := range o {
		err = errors.Join(single.apply(a))
	}

	return err
}

type option func(*config) error

func (o option) apply(cfg *config) error {
	return o(cfg)
}

// ShowBanner is used to show the Echo banner.
func ShowBanner() Option {
	return option(func(c *config) error {
		c.echo.HideBanner = false
		return nil
	})
}

// SetLogger is used to set the logger.
// It overwrites any previous
func SetLogger(logger *slog.Logger) Option {
	return option(func(c *config) error {
		if logger == nil {
			return ErrInvalidLogger
		}
		c.logger = logger
		return nil
	})
}

// SetSecrets
func SetSecrets(secrets ...string) Option {
	return option(func(c *config) error {

		if len(secrets) < 1 {
			return errors.New("must have at least 1 argument")
		}

		c.secrets = secrets

		return nil
	})
}

// WithMiddleware adds a middleware to the base
func WithMiddleware(middleware ...echo.MiddlewareFunc) Option {
	return option(func(c *config) error {
		c.echo.Use(middleware...)
		return nil
	})
}

// WithRateLimit adds a rate limit middleware to the base
func WithRateLimit(limit rate.Limit) Option {
	return WithMiddleware(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(limit)))
}

// SPA is a shorthand for Single Page Application
func SPA(path string, fs embed.FS) Option {
	return SinglePageApplication(path, fs)
}

// SinglePageApplication is used to server a Single Page Application from an embed.FS
func SinglePageApplication(path string, fs embed.FS) Option {
	return WithMiddleware(middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:       path,
			HTML5:      true,
			Filesystem: http.FS(fs),
		},
	))
}

// Assets is used to serve embeded assets, i.e., css, js, or img.
func Assets(path string, fs embed.FS, browse bool) Option {
	return WithMiddleware(middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:       path,
			Browse:     browse,
			Filesystem: http.FS(fs),
		},
	))
}

// Path
func Path(method string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) Option {
	return option(func(c *config) error {

		if !slices.Contains([]string{
			http.MethodConnect,
			http.MethodDelete,
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
			http.MethodPatch,
			http.MethodPost,
			http.MethodPut,
			http.MethodTrace,
		}, method) {
			return ErrInvalidHTTPMethod
		}

		c.echo.Add(method, path, handler, middleware...)

		return nil
	})
}

// GraphQL
func GraphQL(path string, websockets bool, input *handler.Server, middleware ...echo.MiddlewareFunc) Option {

	fn := func(c echo.Context) error {
		input.ServeHTTP(c.Response(), c.Request())
		return nil
	}

	list := Options{
		Path(http.MethodPost, path, fn, middleware...),
	}

	if websockets {
		list = append(list, Path(http.MethodGet, path, fn, middleware...))
	}

	return list
}

// Playground adds the GraphQL Playground. This is usually for debugging.
func Playground(path string, middleware ...echo.MiddlewareFunc) Option {
	return Path(http.MethodGet, path, func(c echo.Context) error {
		playground.Handler("GraphQL", path).ServeHTTP(c.Response(), c.Request())
		return nil
	}, middleware...)
}

// Timeout sets the timeout to all routes.
func Timeout(duration time.Duration) Option {

	if duration < 0 {
		return option(func(c *config) error { return ErrInvalidDuration })
	}

	return WithMiddleware(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: duration,
	}))
}

// Quiet hides the port information and the banner.
func Quiet() Option {
	return option(func(c *config) error {
		c.echo.HideBanner = true
		c.echo.HidePort = true
		return nil
	})
}

func Renderer(r echo.Renderer) Option {
	return option(func(c *config) error {

		if r == nil {
			return ErrEmptyRenderer
		}

		c.echo.Renderer = r

		return nil
	})
}

func New(opts ...Option) (*echo.Echo, error) {

	// Start Echo
	e := echo.New()

	e.HideBanner = true

	// Start the config.
	c := config{
		logger:  slog.Default(),
		session: scs.New(),
		echo:    e,
	}

	// Prepare an error variable for use with the user provided options.
	var err error

	// Overwrite the config with the user provided options.
	for _, opt := range opts {
		err = errors.Join(opt.apply(&c))
	}

	// Check if there are any errors from the user provided options.
	if err != nil {
		return nil, err
	}

	// Use the desired logger.
	c.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(ctx echo.Context, v middleware.RequestLoggerValues) error {

			attrs := []slog.Attr{
				slog.String("id", v.RequestID),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("host", v.Host),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("user_agent", v.UserAgent),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("content_length", v.ContentLength),
				slog.Int64("response_size", v.ResponseSize),
			}

			if v.Error != nil {
				attrs = append(attrs, slog.String("error", v.Error.Error()))
			}

			c.logger.LogAttrs(ctx.Request().Context(), slog.LevelInfo, "", attrs...)

			return nil
		},
		HandleError:      true,
		LogLatency:       true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogRequestID:     true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
	}))

	return c.echo, nil
}
