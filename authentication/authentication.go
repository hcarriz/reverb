package authentication

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/sessions"
	"github.com/hcarriz/reverb/authentication/provider"
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
	ErrSessionsIsNil      = errors.New("provided Sessions is nil")
)

const (
	DefaultProvider = "provider"
)

type Session interface {
	GetString(ctx context.Context, key string) string
	Token(context.Context) string
	Destroy(context.Context) error
	Put(ctx context.Context, key string, val any)
	PopBool(ctx context.Context, key string) bool
	RenewToken(context.Context) error
	Commit(context.Context) (token string, expires time.Time, err error)
}

// Log interface that this package uses. By default it is `log/slog`
type Log interface {
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

// gothID is for the id of the user from the provider.
// userID is for the id of the user from the database.
type DB interface {
	GetUserWithSession(ctx context.Context, userID string, token string) (user any, err error)               // Get the user with the userID and the active token.
	GetUserID(ctx context.Context, gothID string) (userID string, err error)                                 // GetUserID takes a id from goth and searches for a user in the database. It returns the id of the user
	GetUserIDFromToken(ctx context.Context, apiToken string) (string, error)                                 // Get the user id from the api key
	UpdateUserInfo(ctx context.Context, id string, email string, name string) error                          // Update a user with the given id from the database with the given email and name.
	CreateOrUpdateUser(ctx context.Context, gothID, provider, email, name string) (userID string, err error) // Create or update a user for the given gothID
	UserDisabled(ctx context.Context, userID string) (disabled bool, err error)                              // Check if a user is disabled.
	AddSessionToUser(ctx context.Context, gothID string, session string) error                               // Adds the session to the user with the connected goth id. Returns an error if a user is not connected to the goth id.
}

type auth struct {
	session   Session
	backend   *url.URL
	frontend  *url.URL
	logger    Log
	paths     paths
	names     names
	db        DB
	providers []provider.Provider
}

type paths struct {
	afterError  string
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

func SetSessions(s Session) Option {
	return option(func(a *auth) error {
		if s == nil {
			return ErrSessionsIsNil
		}

		a.session = s

		return nil
	})
}

// WithProvider adds a provider. Source is required for Okta, Nextcloud, and OpenID Providers.
func WithProvider(p provider.Provider, key, secret, callbackDomain, source string) Option {
	return option(func(a *auth) error {

		if !provider.Validate(p) {
			return ErrInvalidProvider
		}

		if slices.Contains([]provider.Provider{provider.Okta, provider.NextCloud, provider.OpenID}, p) {
			if _, err := url.Parse(source); err != nil {
				return err
			}
		}

		u, err := url.Parse(callbackDomain)
		if err != nil {
			return err
		}

		u.Path = fmt.Sprintf("/auth/callback/%s", p)

		d := cloneURL(u)
		d.Path = ""

		if err := p.Use(key, secret, callbackDomain, source); err != nil {
			return err
		}

		a.providers = append(a.providers, p)

		return nil
	})
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

// SetDatabase sets the database that will be used.
func SetDatabase(db DB) Option {
	return option(func(a *auth) error {

		if db == nil {
			return ErrEmptyArgument
		}

		a.db = db

		return nil
	})
}

func SetPaths(login, logout, profile, err string) Option {
	return option(func(a *auth) error {

		a.paths = paths{
			afterError:  err,
			afterLogin:  login,
			afterLogout: logout,
			profile:     profile,
		}

		return nil
	})
}

func SetNames(sesson, provider, refetch, addition string) Option {
	return option(func(a *auth) error {

		a.names = names{
			session:  sesson,
			provider: provider,
			refetch:  refetch,
			addition: addition,
		}

		return nil
	})
}

func New(g *echo.Echo, opts ...Option) error {

	a := &auth{
		logger:  slog.Default(),
		session: scs.New(),
		paths: paths{
			afterError:  "/",
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

	// Sort providers
	slices.SortFunc(a.providers, func(a, b provider.Provider) int {
		switch {
		case a.String() > b.String():
			return 1
		case a.String() < b.String():
			return -1
		default:
			return 0
		}
	})

	// Check for database
	if a.db == nil {
		return ErrMissingDB
	}

	group := g.Group("/auth",
		MiddlewareBearerToken(a.db),
		MiddlewareSessionManager(a.session, a.names.session),
		MiddlewareOIDC(),
	)

	group.GET(fmt.Sprintf("/add/:%s", a.names.provider), a.addExistingAccount, MiddlewareMustBeAuthenticated(a.db))
	group.GET(fmt.Sprintf("/login/:%s", a.names.provider), a.login)
	group.GET(fmt.Sprintf("/callback/:%s", a.names.provider), a.callback)
	group.GET("/logout", a.logout, MiddlewareMustBeAuthenticated(a.db))
	group.GET("/providers", a.listProviders)
	group.GET(fmt.Sprintf("/refetch/:%s", a.names.provider), a.refetch, MiddlewareMustBeAuthenticated(a.db))
	group.GET("/whoami", a.whoami)

	return nil
}

func (a *auth) redirect(c echo.Context, path string, internal bool) error {

	if a.frontend != nil && internal {

		p1 := path

		u2 := cloneURL(a.frontend)
		u2.Path = path

		path = u2.String()

		a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "transforming internal url", slog.String("from", p1), slog.String("to", path))

	}

	a.logger.LogAttrs(context.Background(), slog.LevelDebug, "redirecting", slog.Bool("internal", internal), slog.String("to", path))

	return c.Redirect(http.StatusTemporaryRedirect, path)
}

func (a *auth) err(c echo.Context, msg string, err error, attrs ...slog.Attr) error {
	a.logger.LogAttrs(c.Request().Context(), slog.LevelError, msg, attrs...)
	return c.String(http.StatusMethodNotAllowed, err.Error())
}

func (a *auth) listProviders(c echo.Context) error {

	list := make(map[string]string)

	for _, single := range a.providers {
		list[single.String()] = single.Pretty()
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "displaying available providers", slog.Any("providers", list))

	return c.JSON(http.StatusOK, list)
}

func (a *auth) whoami(c echo.Context) error {

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "who am i", slog.String("session", a.names.session))

	userID := a.session.GetString(c.Request().Context(), a.names.session)
	if userID == "" {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelWarn, "user does not have session")
		return c.NoContent(http.StatusUnauthorized)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "found user session", slog.String("user_id", userID))

	token := a.session.Token(c.Request().Context())

	result, err := a.db.GetUserWithSession(c.Request().Context(), userID, token)
	if err != nil {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelError, "unable to get user from database", slog.String("user_id", userID), slerr(err))
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

		a.logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "signing in user", slog.String("provider", provider))

		redir, err := gothic.GetAuthURL(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
		if err != nil {
			return a.err(c, "unable to get auth url", err)
		}

		a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "redirecting user to identity provider", slog.String("provider", provider), slog.String("url", redir))

		return a.redirect(c, redir, false)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "success, redirecting", slog.String("uri", a.paths.afterLogin), slog.String("user", usr.UserID))

	return a.redirect(c, a.paths.afterLogin, true)

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

	return a.redirect(c, a.paths.afterLogout, true)
}

func (a *auth) addExistingAccount(c echo.Context) error {

	provider := c.Param(a.names.provider)

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "attempting to add existing account", slog.String("provider", provider))

	userID := a.session.GetString(c.Request().Context(), a.names.session)
	if userID != "" {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelError, "no user id in session")
		return c.NoContent(http.StatusUnauthorized)
	}

	a.session.Put(c.Request().Context(), a.names.addition, true)

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "redirecting to provider", slog.String("provider", provider))

	u, err := gothic.GetAuthURL(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	if err != nil {
		return a.err(c, "unable to get auth url", err, slog.String("provider", provider))
	}

	return a.redirect(c, u, false)
}

func (a *auth) callback(c echo.Context) error {

	if a.session.PopBool(c.Request().Context(), a.names.refetch) {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "performing refetch callback")
		return a.callbackRefetch(c)
	}

	if a.session.PopBool(c.Request().Context(), a.names.addition) {
		a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "performing addition callback")
		return a.callbackAddition(c)
	}

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "performing login callback")
	return a.callbackLogin(c)
}

func (a *auth) refetch(c echo.Context) error {

	provider := c.Param(a.names.provider)

	a.logger.LogAttrs(c.Request().Context(), slog.LevelDebug, "refetching data", slog.String("provider", provider))

	usr := a.session.GetString(c.Request().Context(), a.names.session)
	if usr != "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	u, err := gothic.GetAuthURL(c.Response(), gothic.GetContextWithProvider(c.Request(), provider))
	if err != nil {
		return a.err(c, "unable to get auth url", err, slog.String("provider", provider))
	}

	a.session.Put(c.Request().Context(), a.names.refetch, true)

	if _, _, err := a.session.Commit(c.Request().Context()); err != nil {
		return a.err(c, "unable to commit session", err)
	}

	return a.redirect(c, u, false)

}

func gothicName(g goth.User) string {
	return g.Name
}
