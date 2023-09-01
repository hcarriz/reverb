package authentication

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

var (
	ErrEmptyArgument      = errors.New("argument can not be empty")
	ErrInvalidProvider    = errors.New("provider is invalid")
	ErrLoggerIsInvalid    = errors.New("provided logger is nil")
	ErrMissingDB          = errors.New("missing database")
	ErrMissingDomain      = errors.New("missing domain")
	ErrMissingProvider    = errors.New("missing provider")
	ErrNotSetup           = errors.New("this function has not been completed")
	ErrSessionsStoreIsNil = errors.New("provided sessions.Store is nil")
)

const (
	DefaultProvider = "provider"
)

// Log interface that this package uses. By default it is `log/slog`
type Log interface {
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

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
	backend   *url.URL
	frontend  *url.URL
	logger    Log
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

type Options []Option

func (o Options) apply(a *auth) error {

	var err error

	for _, option := range o {
		err = errors.Join(option.apply(a))
	}

	return err
}

type option func(*auth) error

func (o option) apply(a *auth) error {
	return o(a)
}

func SetStore(store sessions.Store) Option {
	return option(func(a *auth) error {

		if store == nil {
			return ErrSessionsStoreIsNil
		}

		gothic.Store = store

		return nil
	})
}

func SetLogger(logger Log) Option {
	return option(func(a *auth) error {

		if logger == nil {
			return ErrLoggerIsInvalid
		}

		a.logger = logger

		return nil
	})
}

// SetFrontend sets the url for the frontend to use.
func SetFrontend(raw string) Option {
	return option(func(a *auth) error {

		if raw == "" {
			return ErrEmptyArgument
		}

		u, err := url.Parse(raw)
		if err != nil {
			return err
		}

		a.frontend = u

		return nil
	})
}

// SetBackend sets the url for the backend to use.
func SetBackend(raw string) Option {
	return option(func(a *auth) error {

		if raw == "" {
			return ErrEmptyArgument
		}

		u, err := url.Parse(raw)
		if err != nil {
			return err
		}

		a.backend = u

		return nil
	})
}

func New(g *echo.Echo, opts ...Option) error {

	a := &auth{
		logger:  slog.Default(),
		session: scs.New(),
		paths: paths{
			afterLogin:  "/",
			afterLogout: "/",
			profile:     "/",
		},
		names: names{
			session:  "user_session",
			provider: DefaultProvider,
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

	group := g.Group("/auth",
		MiddlewareBearerToken(a.db),
		MiddlewareSessionManager(a.session, a.names.session),
		MiddlewareOIDC(),
		MiddlewareSetIPAddress(),
	)

	group.GET(fmt.Sprintf("/add/:%s", a.names.provider), a.addExistingAccount, MiddlewareMustBeAuthenticated(a.db))
	group.GET(fmt.Sprintf("/login/:%s", a.names.provider), a.login)
	group.GET(fmt.Sprintf("/callback/:%s", a.names.provider), a.callback)
	group.GET("/logout", a.logout, MiddlewareMustBeAuthenticated(a.db))
	group.GET("/providers", a.listProviders)
	group.GET("/refetch", a.refetch, MiddlewareMustBeAuthenticated(a.db))
	group.GET("/whoami", a.whoami)

	return nil
}

func (a *auth) redir(c echo.Context, path string) error {

	if a.frontend != nil {

		u2 := cloneURL(a.frontend)
		u2.Path = path

		path = u2.String()

	}

	return c.Redirect(http.StatusFound, path)
}

func (a *auth) err(c echo.Context, err error) error {
	return c.String(http.StatusMethodNotAllowed, err.Error())
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
		a.logger.LogAttrs(c.Request().Context(), slog.LevelWarn, "user does not have session")
		return c.NoContent(http.StatusUnauthorized)
	}

	token := a.session.Token(c.Request().Context())

	result, err := a.db.GetUserWithSession(c.Request().Context(), userID, token)
	if err != nil {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelError, "unable to get user", slog.String("user", userID), slerr(err))
		return c.NoContent(http.StatusUnauthorized)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "found user", slog.String("id", userID), slog.Any("data", result))

	return c.JSON(http.StatusOK, result)
}

func (a *auth) login(c echo.Context) error {

	provider := c.Param(a.names.provider)

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "attempting to login", slog.String("provider", provider))

	usr, err := gothic.CompleteUserAuth(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	if err != nil {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelWarn, "user is not signed in")
		a.logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "signing in user")
		gothic.BeginAuthHandler(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
		return nil
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "success, redirecting", slog.String("uri", a.paths.afterLogin), slog.String("user", usr.UserID))

	return a.redir(c, a.paths.afterLogin)

}

func (a *auth) logout(c echo.Context) error {

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "attempting to logout")

	if provider := a.session.GetString(c.Request().Context(), a.names.provider); provider != "" {
		if usr := a.session.GetString(c.Request().Context(), a.names.session); usr != "" {
			if err := gothic.Logout(c.Response(), gothic.GetContextWithProvider(c.Request(), provider)); err != nil {
				a.logger.LogAttrs(c.Request().Context(), slog.LevelError, "unable to logout", slerr(err))
			}
		}
	}

	if err := a.session.Destroy(c.Request().Context()); err != nil {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelError, "unable to destroy session", slerr(err))
	}

	return a.redir(c, a.paths.afterLogout)
}

func (a *auth) addExistingAccount(c echo.Context) error {

	userID := a.session.GetString(c.Request().Context(), a.names.session)
	if userID != "" {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelError, "no user")
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
			return a.authError(c, "unable to connect user with authentication", err, slog.String("user", userID), slog.String("provider", provider), slog.String("connection_id", connectionID))
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
