package cors

import (
	"errors"
	"net/url"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	ErrHasQuery          = errors.New("includes query")
	ErrMissingScheme     = errors.New("missing scheme")
	ErrHasPath           = errors.New("has path")
	ErrInvalidHeader     = errors.New("invalid headers")
	ErrOriginFuncDefined = errors.New("cors.OriginFunc has been used")
	ErrOriginsDefined    = errors.New("cores.Origins has been users")
	ErrInvalidAge        = errors.New("max age can not be negative")

	AcceptableHeaders = []string{
		echo.HeaderAccept,
		echo.HeaderAcceptEncoding,
		echo.HeaderAllow,
		echo.HeaderAuthorization,
		echo.HeaderContentDisposition,
		echo.HeaderContentEncoding,
		echo.HeaderContentLength,
		echo.HeaderContentType,
		echo.HeaderCookie,
		echo.HeaderSetCookie,
		echo.HeaderIfModifiedSince,
		echo.HeaderLastModified,
		echo.HeaderLocation,
		echo.HeaderRetryAfter,
		echo.HeaderUpgrade,
		echo.HeaderVary,
		echo.HeaderWWWAuthenticate,
		echo.HeaderXForwardedFor,
		echo.HeaderXForwardedProto,
		echo.HeaderXForwardedProtocol,
		echo.HeaderXForwardedSsl,
		echo.HeaderXUrlScheme,
		echo.HeaderXHTTPMethodOverride,
		echo.HeaderXRealIP,
		echo.HeaderXRequestID,
		echo.HeaderXCorrelationID,
		echo.HeaderXRequestedWith,
		echo.HeaderServer,
		echo.HeaderOrigin,
		echo.HeaderCacheControl,
		echo.HeaderConnection,
		echo.HeaderAccessControlRequestMethod,
		echo.HeaderAccessControlRequestHeaders,
		echo.HeaderAccessControlAllowOrigin,
		echo.HeaderAccessControlAllowMethods,
		echo.HeaderAccessControlAllowHeaders,
		echo.HeaderAccessControlAllowCredentials,
		echo.HeaderAccessControlExposeHeaders,
		echo.HeaderAccessControlMaxAge,
		echo.HeaderStrictTransportSecurity,
		echo.HeaderXContentTypeOptions,
		echo.HeaderXXSSProtection,
		echo.HeaderXFrameOptions,
		echo.HeaderContentSecurityPolicy,
		echo.HeaderContentSecurityPolicyReportOnly,
		echo.HeaderXCSRFToken,
		echo.HeaderReferrerPolicy,
	}
)

type Option interface {
	apply(*middleware.CORSConfig) error
}

type option func(*middleware.CORSConfig) error

func (o option) apply(c *middleware.CORSConfig) error {
	return o(c)
}

func Skipper(skipper middleware.Skipper) Option {
	return option(func(c *middleware.CORSConfig) error {
		c.Skipper = skipper
		return nil
	})
}

func Origins(list ...string) Option {
	return option(func(c *middleware.CORSConfig) error {

		if c.AllowOriginFunc != nil {
			return ErrOriginFuncDefined
		}

		for _, single := range list {

			if single == "*" || single == "?" {
				continue
			}

			u, err := url.Parse(single)
			if err != nil {
				return err
			}

			// TODO: Add more checks.
			switch {
			case u.Query().Encode() != "":
				return ErrHasQuery
			case u.Path != "":
				return ErrHasPath
			case u.Scheme == "":
				return ErrMissingScheme
			}

		}

		c.AllowOrigins = append(c.AllowOrigins, list...)

		return nil
	})
}

// Headers is used to give the headers that are allowed, others are not allowed.
func Headers(list ...string) Option {
	return option(func(c *middleware.CORSConfig) error {

		for _, single := range list {
			if !slices.Contains(AcceptableHeaders, single) {
				return ErrInvalidHeader
			}
		}

		c.AllowHeaders = append(c.AllowHeaders, list...)

		return nil
	})
}

// OriginFunc allows for a function that checks the origin to be used by the CORS middleware.
func OriginFunc(f func(origin string) (bool, error)) Option {
	return option(func(c *middleware.CORSConfig) error {

		if len(c.AllowOrigins) > 0 {
			return ErrOriginsDefined
		}

		c.AllowOriginFunc = f
		return nil
	})
}

// Methods allows for certain methods to be allowed by the CORS middleware.
func Methods(methods ...string) Option {
	return option(func(c *middleware.CORSConfig) error {

		//TODO: check the methods to see if they're valid

		c.AllowMethods = append(c.AllowMethods, methods...)

		return nil
	})
}

// Credentials allows for credentials to be used.
func Credentials() Option {
	return option(func(c *middleware.CORSConfig) error {
		c.AllowCredentials = true
		return nil
	})
}

func MaxAge(age int) Option {
	return option(func(c *middleware.CORSConfig) error {
		if age < 0 {
			return ErrInvalidAge
		}

		c.MaxAge = age

		return nil
	})
}

func New(opts ...Option) (echo.MiddlewareFunc, error) {

	config := middleware.CORSConfig{}

	var err error

	for _, opt := range opts {
		err = errors.Join(opt.apply(&config))
	}

	if err != nil {
		return nil, err
	}

	return middleware.CORSWithConfig(config), nil
}
