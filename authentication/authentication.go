package authentication

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

var (
	ErrMissingDB = errors.New("missing database")
)

type DB interface {
	GetUserWithSession(ctx context.Context, userID string, session string) (any, error)
	GetUser(ctx context.Context, userID string) (any, error)
	GetUserIDFromToken(ctx context.Context, token string) (string, error)
	UpdateUserInfo(ctx context.Context, id string, email string, name string) error
	CreateConnection(ctx context.Context, oidcUserID string, provider string) (string, error)
	UserConnection(ctx context.Context, userID string, connectionID string) error
}

type auth struct {
	session   *scs.SessionManager
	logger    *slog.Logger
	paths     paths
	names     names
	db        DB
	providers []Provider
}

type paths struct {
	whoami      string
	afterLogin  string
	afterLogout string
	profile     string
}

type routes struct {
	whoami string
	logout string
}

type names struct {
	session  string
	provider string
	refetch  string
	addition string
}

type Option interface {
	apply(*auth) error
}

func New(g *echo.Group, opts ...Option) error {

	a := &auth{
		logger: slog.Default(),
		paths: paths{
			whoami:      "/auth/whoami",
			afterLogin:  "/",
			afterLogout: "/",
			profile:     "/profile",
		},
		names: names{
			session:  "user_session",
			provider: "provider",
			refetch:  "user_refetch",
			addition: "add_existing_account",
		},
	}

	var err error

	for _, opt := range opts {
		err = errors.Join(opt.apply(a))
	}

	if err != nil {
		return err
	}

	// Check for database
	if a.db == nil {
		return ErrMissingDB
	}

	g.GET(fmt.Sprintf("/add/:%s", a.names.provider), a.addExistingAccount)
	g.GET(fmt.Sprintf("/login/:%s", a.names.provider), a.login)
	g.GET(fmt.Sprintf("/callback/:%s", a.names.provider), a.callback)
	g.GET(fmt.Sprintf("/logout/:%s", a.names.provider), a.logout)
	g.GET("/providers", a.listProviders)
	g.GET("/refetch", a.refetch)
	g.GET("/whoami", a.whoami)

	return nil
}

func (a *auth) listProviders(c echo.Context) error {

	// Completed

	list := []string{}

	for _, single := range a.providers {
		list = append(list, single.display)
	}

	return c.JSON(http.StatusOK, list)
}

func (a *auth) whoami(c echo.Context) error {

	userID := a.session.GetString(c.Request().Context(), a.names.session)
	if userID == "" {
		a.logger.Warn("user does not have session")
		return c.NoContent(http.StatusUnauthorized)
	}

	token := a.session.Token(c.Request().Context())

	result, err := a.db.GetUserWithSession(c.Request().Context(), userID, token)
	if err != nil {
		a.logger.Error("unable to get user", slog.String("user", userID), slog.String("error", err.Error()))
		return c.NoContent(http.StatusUnauthorized)
	}

	a.logger.Debug("found user", slog.String("id", userID), slog.Any("data", result))

	return c.JSON(http.StatusOK, result)
}

func (a *auth) login(c echo.Context) error {

	provider := c.Param(a.names.provider)

	if _, err := gothic.CompleteUserAuth(c.Response(), gothic.GetContextWithProvider(c.Request(), provider)); err != nil {
		gothic.BeginAuthHandler(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
		return nil
	}

	return c.Redirect(http.StatusFound, a.paths.afterLogin)

}

func (a *auth) logout(c echo.Context) error {

	provider := c.Param(a.names.provider)

	usr := a.session.GetString(c.Request().Context(), a.names.session)
	if usr != "" {
		gothic.Logout(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	}

	a.session.Destroy(c.Request().Context())

	return c.Redirect(http.StatusFound, a.paths.afterLogout)
}

// addExistingAccount should be locked down with
func (a *auth) addExistingAccount(c echo.Context) error {

	userID := a.session.GetString(c.Request().Context(), a.names.session)
	if userID != "" {
		a.logger.Error("no user")
		return c.NoContent(http.StatusUnauthorized)
	}

	a.session.Put(c.Request().Context(), a.names.addition, true)

	gothic.BeginAuthHandler(c.Response(), gothic.GetContextWithProvider(c.Request(), c.Param(a.names.provider)))

	return nil

}

func (a *auth) callback(c echo.Context) error {

	provider := c.Param(a.names.provider)

	if err := a.session.RenewToken(c.Request().Context()); err != nil {
		return a.authError(c, "unable to renew token", err)
	}

	u, err := gothic.CompleteUserAuth(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	if err != nil {
		return a.authError(c, "unable to complete", err)
	}

	// a.db.CreateConnection()

	if a.session.GetBool(c.Request().Context(), a.names.refetch) {

		a.db.UpdateUserInfo(c.Request().Context(), u.UserID, u.Email, gothicName(u))

		a.session.Put(c.Request().Context(), a.names.refetch, false)

		return c.Redirect(http.StatusFound, a.paths.profile)

	}

	if a.session.GetBool(c.Request().Context(), a.names.addition) {

		userID := a.session.GetString(c.Request().Context(), a.names.session)
		if userID != "" {
			return a.authError(c, "unauthorized", errors.New("unauthorized"))
		}

		connectionID, err := a.db.CreateConnection(c.Request().Context(), u.UserID, provider)
		if err != nil {
			return a.authError(c, "unable to create connection", err, slog.String("user", userID), slog.String("provider", provider))
		}

		if err := a.db.UserConnection(c.Request().Context(), userID, connectionID); err != nil {
			return a.authError(c, "unable to connect user with authentication", err, slog.String("user", userID), slog.String("provider", provider), slog.String("conenction_id", connectionID))
		}

		return c.Redirect(http.StatusFound, a.paths.profile)

	}

	return c.Redirect(http.StatusFound, a.paths.profile)
}

func (a *auth) refetch(c echo.Context) error {

	provider := c.Param(a.names.provider)

	usr := a.session.GetString(c.Request().Context(), a.names.session)
	if usr != "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	a.session.Put(c.Request().Context(), a.names.refetch, true)

	gothic.BeginAuthHandler(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))

	return nil

}

func (a *auth) authError(c echo.Context, note string, err error, attr ...slog.Attr) error {

	a.logger.LogAttrs(c.Request().Context(), slog.LevelError, note, attr...)

	c.SetCookie(&http.Cookie{
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Name:     a.names.session,
		Path:     "/",
		Value:    "",
	})

	return c.Redirect(http.StatusFound, a.paths.afterLogout)
}

func gothicName(g goth.User) string {
	return ""
}
