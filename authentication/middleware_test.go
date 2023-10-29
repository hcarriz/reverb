package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/hcarriz/reverb/authentication/dummy"
	"github.com/hcarriz/reverb/authentication/provider"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/neilotoole/slogt"
	"github.com/oauth2-proxy/mockoidc"
	session "github.com/spazzymoto/echo-scs-session"
	"github.com/stretchr/testify/require"
)

type MockUser struct {
	Subject           string
	Email             string
	EmailVerified     bool
	PreferredUsername string
	Phone             string
	Address           string
	Groups            []string
}

type mockUserinfo struct {
	Email             string   `json:"email,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	Phone             string   `json:"phone_number,omitempty"`
	Address           string   `json:"address,omitempty"`
	Groups            []string `json:"groups,omitempty"`
}

func (u *MockUser) Userinfo(scope []string) ([]byte, error) {
	user := u.scopedClone(scope)

	info := &mockUserinfo{
		Email:             user.Email,
		PreferredUsername: user.PreferredUsername,
		Phone:             user.Phone,
		Address:           user.Address,
		Groups:            user.Groups,
	}

	return json.Marshal(info)
}

func (mu *MockUser) ID() string {
	return mu.Subject
}

type mockClaims struct {
	*mockoidc.IDTokenClaims
	Email             string   `json:"email,omitempty"`
	EmailVerified     bool     `json:"email_verified,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	Phone             string   `json:"phone_number,omitempty"`
	Address           string   `json:"address,omitempty"`
	Groups            []string `json:"groups,omitempty"`
}

func (u *MockUser) Claims(scope []string, claims *mockoidc.IDTokenClaims) (jwt.Claims, error) {
	user := u.scopedClone(scope)

	return &mockClaims{
		IDTokenClaims:     claims,
		Email:             user.Email,
		EmailVerified:     user.EmailVerified,
		PreferredUsername: user.PreferredUsername,
		Phone:             user.Phone,
		Address:           user.Address,
		Groups:            user.Groups,
	}, nil
}

func (u *MockUser) scopedClone(scopes []string) *MockUser {
	clone := &MockUser{
		Subject: u.Subject,
	}
	for _, scope := range scopes {
		switch scope {
		case "profile":
			clone.PreferredUsername = u.PreferredUsername
			clone.Address = u.Address
			clone.Phone = u.Phone
		case "email":
			clone.Email = u.Email
			clone.EmailVerified = u.EmailVerified
		case "groups":
			clone.Groups = append(make([]string, 0, len(u.Groups)), u.Groups...)
		}
	}
	return clone
}

func TestSetUser(t *testing.T) {

	check := require.New(t)
	sm := scs.New()
	sl := slogt.New(t)
	db := &dummy.DB{}
	a := &auth{}

	// http client that doesn't redirect
	cl := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	m, err := mockoidc.Run()
	check.NoError(err)

	m.QueueUser(&MockUser{})

	defer m.Shutdown()

	gothic.Store = sessions.NewCookieStore([]byte("insecure_key"))

	oidcConfig := m.Config()

	opts := Options{
		SetDatabase(db),
		SetSessions(sm),
		SetLogger(sl),
		SetPaths("/", "/", "/", "/"),
		SetNames("app_session", "provider", "app_refetch", "app_addition"),
		WithProvider(provider.OpenID, oidcConfig.ClientID, oidcConfig.ClientSecret, "", fmt.Sprintf("%s/.well-known/openid-configuration", oidcConfig.Issuer)),
	}

	for _, opt := range opts {
		check.NoError(opt.apply(a))
	}

	e := echo.New()

	mw := []echo.MiddlewareFunc{
		session.LoadAndSave(sm),
		// MiddlewareBearerToken(db),
		// MiddlewareSessionManager(sm, "app_session"),
		// a.middlewaresm(),
		// MiddlewareOIDC(),
	}

	lgn := wares(a.login, mw...)

	clb := wares(a.callback, mw...)

	// Test login

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	ctx.SetPath("/auth/login/:provider")
	ctx.SetParamNames("provider")
	ctx.SetParamValues("openid-connect")

	check.NoError(lgn(ctx))

	check.NotNil(rec.Result())

	check.Equal(http.StatusTemporaryRedirect, rec.Result().StatusCode)

	cookies := rec.Result().Cookies()

	for _, cookie := range cookies {
		t.Log(cookie.Name)
	}

	goingto := rec.Result().Header.Get("Location")
	check.NotEmpty(goingto)

	req, err = http.NewRequestWithContext(ctx.Request().Context(), http.MethodGet, goingto, nil)
	check.NoError(err)

	resp, err := cl.Do(req)
	check.NoError(err)
	check.Equal(http.StatusFound, resp.StatusCode)

	callback := resp.Header.Get("Location")
	check.NotEmpty(callback)

	t.Log(callback)

	req = httptest.NewRequest(http.MethodGet, callback, nil)
	rec = httptest.NewRecorder()

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	ctx = e.NewContext(req, rec)

	ctx.SetPath("/auth/callback/:provider")
	ctx.SetParamNames("provider")
	ctx.SetParamValues("openid-connect")

	check.NoError(clb(ctx))

	check.Equal(http.StatusTemporaryRedirect, rec.Result().StatusCode)
}
