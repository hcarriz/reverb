package slocho

import (
	"context"
	"errors"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	ErrSkipperNil  = errors.New("skipper is nil")
	ErrLoggerIsNil = errors.New("logger is nil")
)

type config struct {
	msg                  string
	level                slog.Level
	logger               Log
	skipper              middleware.Skipper
	beforeFunc           func(echo.Context)
	excludeErr           bool
	excludeLatency       bool
	excludeProtocol      bool
	excludeRemoteIP      bool
	excludeHost          bool
	excludeMethod        bool
	excludeURI           bool
	excludeURIPath       bool
	excludeRoutePath     bool
	excludeRequestID     bool
	excludeReferer       bool
	excludeUserAgent     bool
	excludeStatus        bool
	excludeError         bool
	excludeContentLength bool
	excludeResponseSize  bool
	headers              []string
	queryParams          []string
	formValues           []string
}

type Log interface {
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

type Option interface {
	apply(*config) error
}

type option func(*config) error

func (o option) apply(c *config) error {
	return o(c)
}

// Skipper lets you use the skipper of your choice.
func Skipper(sk middleware.Skipper) Option {
	return option(func(c *config) error {

		if sk == nil {
			return ErrSkipperNil
		}

		c.skipper = sk

		return nil
	})
}

// Logger sets the logger to be used.
func Logger(logger Log) Option {
	return option(func(c *config) error {

		if logger == nil {
			return ErrLoggerIsNil
		}

		c.logger = logger

		return nil
	})
}

// Message sets the message that is shown for every log entry.
func Message(msg string) Option {
	return option(func(c *config) error {
		c.msg = msg
		return nil
	})
}

// Level sets the level that the logger will write at.
func Level(lvl slog.Level) Option {
	return option(func(c *config) error {
		c.level = lvl
		return nil
	})
}

// WithHeaders lets you choose which headers you want added to the log entry.
func WithHeaders(list ...string) Option {
	return option(func(c *config) error {

		c.headers = append(c.headers, list...)

		return nil
	})
}

// WithQueryParams lets you choose which query params you want added to the log entry.
func WithQueryParams(list ...string) Option {
	return option(func(c *config) error {
		c.queryParams = append(c.queryParams, list...)
		return nil
	})
}

// WithFormValues lets you choose which form values you want added to the log entry.
func WithFormValues(list ...string) Option {
	return option(func(c *config) error {
		c.formValues = append(c.formValues, list...)
		return nil
	})
}

func New(opts ...Option) (echo.MiddlewareFunc, error) {

	var (
		err error
		c   = config{
			logger: slog.Default(),
			msg:    "request",
			level:  slog.LevelInfo,
		}
	)

	for _, opt := range opts {
		err = errors.Join(err, opt.apply(&c))
	}

	if err != nil {
		return nil, err
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper:        c.skipper,
		BeforeNextFunc: c.beforeFunc,
		LogValuesFunc: func(ctx echo.Context, v middleware.RequestLoggerValues) error {

			attrs := []slog.Attr{}

			if !c.excludeLatency {
				attrs = append(attrs, slog.Duration("latency", v.Latency))
			}

			if !c.excludeProtocol {
				attrs = append(attrs, slog.String("protocol", v.Protocol))
			}

			if !c.excludeRemoteIP {
				attrs = append(attrs, slog.String("remote_ip", v.RemoteIP))
			}

			if !c.excludeHost {
				attrs = append(attrs, slog.String("host", v.Host))
			}

			if !c.excludeMethod {
				attrs = append(attrs, slog.String("method", v.Method))
			}

			if !c.excludeURI {
				attrs = append(attrs, slog.String("uri", v.URI))
			}

			if !c.excludeURIPath {
				attrs = append(attrs, slog.String("uri_path", v.URIPath))
			}

			if !c.excludeRoutePath {
				attrs = append(attrs, slog.String("route_path", v.RoutePath))
			}

			if !c.excludeRequestID {
				attrs = append(attrs, slog.String("id", v.RequestID))
			}

			if !c.excludeReferer {
				attrs = append(attrs, slog.String("referrer", v.Referer))
			}

			if !c.excludeUserAgent {
				attrs = append(attrs, slog.String("user_agent", v.UserAgent))
			}

			if !c.excludeStatus {
				attrs = append(attrs, slog.Int("status", v.Status))
			}

			if !c.excludeContentLength {
				attrs = append(attrs, slog.String("content_length", v.ContentLength))
			}

			if !c.excludeResponseSize {
				attrs = append(attrs, slog.Int64("response_size", v.ResponseSize))
			}

			if len(c.headers) > 0 {
				attrs = append(attrs, convert("headers", v.Headers)...)
			}

			if len(c.queryParams) > 0 {
				attrs = append(attrs, convert("query_params", v.QueryParams)...)
			}

			if len(c.formValues) > 0 {
				attrs = append(attrs, convert("form_values", v.FormValues)...)
			}

			if v.Error != nil && !c.excludeErr {
				attrs = append(attrs, slog.String("error", v.Error.Error()))
			}

			c.logger.LogAttrs(ctx.Request().Context(), c.level, c.msg, attrs...)

			return nil
		},
		HandleError:      !c.excludeErr,
		LogLatency:       !c.excludeLatency,
		LogProtocol:      !c.excludeProtocol,
		LogRemoteIP:      !c.excludeRemoteIP,
		LogHost:          !c.excludeHost,
		LogMethod:        !c.excludeMethod,
		LogURI:           !c.excludeURI,
		LogURIPath:       !c.excludeURIPath,
		LogRoutePath:     !c.excludeRoutePath,
		LogRequestID:     !c.excludeRequestID,
		LogReferer:       !c.excludeReferer,
		LogUserAgent:     !c.excludeUserAgent,
		LogStatus:        !c.excludeStatus,
		LogError:         !c.excludeError,
		LogContentLength: !c.excludeContentLength,
		LogResponseSize:  !c.excludeResponseSize,
		LogHeaders:       c.headers,
		LogQueryParams:   c.queryParams,
		LogFormValues:    c.formValues,
	}), nil
}

func convert(key string, values map[string][]string) []slog.Attr {

	attrs := []slog.Attr{}
	for k, v := range values {
		g := slog.Group(key, slog.Any(k, v))
		attrs = append(attrs, g)
	}

	return attrs
}
