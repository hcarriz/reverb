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

func (a *auth) MiddlewareSetIPAddress() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ctx := c.Request().Context()

			ctx = context.WithValue(ctx, "ip_address", c.RealIP())

			req := c.Request().Clone(ctx)

			return next(c.Echo().NewContext(req, c.Response()))

		}
	}
}

func (a *auth) MiddlewareSessionManager() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if usrID := a.session.GetString(c.Request().Context(), a.names.session); usrID != "" {
				c = setUser(c, usrID)
			}

			return next(c)
		}
	}
}

func (a *auth) MiddlewareBearerToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if tk := c.Request().Header.Get("Authorization"); tk != "" {

				b := strings.Split(tk, " ")
				if len(b) == 2 {
					tk = b[1]
				}

				usrID, err := a.db.GetUserIDFromToken(c.Request().Context(), tk)
				if err != nil {
					return c.NoContent(http.StatusUnauthorized)
				}

				c = setUser(c, usrID)

			}

			return next(c)
		}
	}
}

func (a *auth) MiddlewareOIDC() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if user, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
				c = setUser(c, user.UserID)
			}

			return next(c)
		}
	}
}
