package authentication

import (
	"context"
	"net/http"
	"strings"

	"ariga.io/sqlcomment"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func setUser(c echo.Context, id string) echo.Context {

	if id == "" {
		return c
	}

	ctx := context.WithValue(c.Request().Context(), "user", id)
	ctx = sqlcomment.WithTag(ctx, "user", id)

	return c.Echo().NewContext(c.Request().Clone(ctx), c.Response())

}

func MiddlewareMustBeAuthenticated(db DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userID, ok := c.Request().Context().Value("user").(string)
			if !ok {
				return c.NoContent(http.StatusUnauthorized)
			}

			if _, err := db.UserDisabled(c.Request().Context(), userID); err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}

func MiddlewareSetIPAddress() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ctx := c.Request().Context()

			ctx = context.WithValue(ctx, "ip_address", c.RealIP())

			req := c.Request().Clone(ctx)

			return next(c.Echo().NewContext(req, c.Response()))

		}
	}
}

func MiddlewareSessionManager(session Session, key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if usrID := session.GetString(c.Request().Context(), key); usrID != "" {
				c = setUser(c, usrID)
			}

			return next(c)
		}
	}
}

func MiddlewareBearerToken(db DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if tk := c.Request().Header.Get("Authorization"); tk != "" {

				b := strings.Split(tk, " ")
				if len(b) == 2 {
					tk = b[1]
				}

				usrID, err := db.GetUserIDFromToken(c.Request().Context(), tk)
				if err != nil {
					return c.NoContent(http.StatusUnauthorized)
				}

				c = setUser(c, usrID)

			}

			return next(c)
		}
	}
}

func MiddlewareOIDC() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if user, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
				c = setUser(c, user.UserID)
			}

			return next(c)
		}
	}
}
