package authentication

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func (a *auth) callbackAddition(c echo.Context) error {

	provider := c.Param(a.names.provider)

	if err := a.session.RenewToken(c.Request().Context()); err != nil {
		return a.callbackError(c, "unable to renew token", err)
	}

	u, err := gothic.CompleteUserAuth(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	if err != nil {
		return a.callbackError(c, "unable to complete user authentication", err)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "user was provided", slog.String("provider", provider), slog.String("user_id", u.UserID))

	return a.callbackError(c, "callback refetch not finished", ErrEmptyArgument)

}
func (a *auth) callbackRefetch(c echo.Context) error {

	provider := c.Param(a.names.provider)

	if err := a.session.RenewToken(c.Request().Context()); err != nil {
		return a.callbackError(c, "unable to renew token", err)
	}

	u, err := gothic.CompleteUserAuth(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	if err != nil {
		return a.callbackError(c, "unable to complete user authentication", err)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "user was provided", slog.String("provider", provider), slog.String("user_id", u.UserID))

	return a.callbackError(c, "callback refetch not finished", ErrEmptyArgument)

}
func (a *auth) callbackLogin(c echo.Context) error {

	provider := c.Param(a.names.provider)

	if err := a.session.RenewToken(c.Request().Context()); err != nil {
		return a.callbackError(c, "unable to renew token", err)
	}

	u, err := gothic.CompleteUserAuth(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	if err != nil {
		return a.callbackError(c, "unable to complete user authentication", err)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "user was provided", slog.String("provider", provider), slog.String("user_id", u.UserID))

	return a.callbackError(c, "callback refetch not finished", ErrEmptyArgument)

}

func (a *auth) callbackError(c echo.Context, msg string, err error, attr ...slog.Attr) error {

	list := []slog.Attr{
		slerr(err),
	}

	list = append(list, attr...)

	a.logger.LogAttrs(c.Request().Context(), slog.LevelError, msg, list...)

	names := []string{
		a.names.provider,
		a.names.refetch,
		a.names.session,
	}

	for _, name := range names {

		c.SetCookie(&http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
		})
	}

	return a.redir(c, a.paths.afterLogout)
}
