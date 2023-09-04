package authentication

import (
	"errors"
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

	id, err := a.db.CreateOrUpdateUser(c.Request().Context(), u.UserID, provider)
	if err != nil {
		return a.callbackError(c, "unable to add user to database", err)
	}

	if id == "" {
		return a.callbackError(c, "unable to get id from database", errors.New("empty id"))
	}

	if ok, err := a.db.UserDisabled(c.Request().Context(), u.UserID); err != nil || ok {

		msg := "unable to check user status"

		if err == nil {
			msg = "user is disabled"
			err = errors.New(msg)

		}

		return a.callbackError(c, msg, err)
	}

	a.session.Put(c.Request().Context(), a.names.session, id)
	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "added user to session", slog.String("user", id))

	token, _, err := a.session.Commit(c.Request().Context())
	if err != nil {
		return a.callbackError(c, "unable to commit to session", err)
	}

	if err := a.db.AddSessionToUser(c.Request().Context(), u.UserID, token); err != nil {
		return a.callbackError(c, "unable to add session to user", err)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "user has been authenticated by identity provider", slog.String("provider", provider), slog.String("url", c.Request().URL.String()))

	return a.redir(c, a.paths.afterLogin)

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

	return c.String(http.StatusMethodNotAllowed, err.Error())
}
