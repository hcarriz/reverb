package authentication

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/sessions"
	"github.com/hcarriz/reverb/authentication/dummy"
	"github.com/hcarriz/reverb/authentication/provider"
	"github.com/icrowley/fake"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/faux"
	"github.com/neilotoole/slogt"
	session "github.com/spazzymoto/echo-scs-session"
	"github.com/stretchr/testify/require"
)

func gzipString(value string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return "err"
	}
	if err := gz.Flush(); err != nil {
		return "err"
	}
	if err := gz.Close(); err != nil {
		return "err"
	}

	return b.String()
}

func ungzipString(value string) string {
	rdata := strings.NewReader(value)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "err"
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return "err"
	}

	return string(s)
}

func wares(f echo.HandlerFunc, mw ...echo.MiddlewareFunc) echo.HandlerFunc {

	for _, s := range mw {
		f = s(f)
	}

	return f
}

type mapKey struct {
	r *http.Request
	n string
}

type ProviderStore struct {
	Store map[mapKey]*sessions.Session
}

func NewProviderStore() *ProviderStore {
	return &ProviderStore{map[mapKey]*sessions.Session{}}
}

func (p ProviderStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	s := p.Store[mapKey{r, name}]
	if s == nil {
		s, err := p.New(r, name)
		return s, err
	}
	return s, nil
}

func (p ProviderStore) New(r *http.Request, name string) (*sessions.Session, error) {
	s := sessions.NewSession(p, name)
	s.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 30,
	}
	p.Store[mapKey{r, name}] = s
	return s, nil
}

func (p ProviderStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	p.Store[mapKey{r, s.Name()}] = s
	return nil
}

func TestWhoAmI(t *testing.T) {

	check := require.New(t)

	store := NewProviderStore()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	gothic.Store = sessions.NewCookieStore([]byte("insecure_key"))

	sm := scs.New()

	e := echo.New()

	c := e.NewContext(req, rec)

	a := &auth{}

	opts := Options{
		SetDatabase(&dummy.DB{}),
		SetLogger(slogt.New(t)),
		SetSessions(sm),
		SetPaths("/", "/", "/", "/"),
		SetNames("app_session", "provider", "app_refetch", "app_addition"),
	}

	for _, opt := range opts {
		check.NoError(opt.apply(a))
	}

	prv := &faux.Provider{}

	goth.ClearProviders()
	check.Len(goth.GetProviders(), 0)
	goth.UseProviders(prv)
	check.Len(goth.GetProviders(), 1)

	h := session.LoadAndSave(sm)

	def := []echo.MiddlewareFunc{
		h,
	}

	wai := wares(a.whoami, def...)
	li := wares(a.login, def...)
	cb := wares(a.callback, def...)
	refe := wares(a.refetch, def...)

	// Who Am I, someone who's not logged in.
	check.NoError(wai(c))

	check.Equal(http.StatusUnauthorized, rec.Result().StatusCode)

	data, err := io.ReadAll(rec.Result().Body)
	check.NoError(err)
	check.Empty(data)

	check.NoError(rec.Result().Body.Close())

	// Login

	req = httptest.NewRequest(http.MethodGet, "/auth/login/faux", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	c.SetPath("/auth/login/:provider")
	c.SetParamNames("provider")
	c.SetParamValues("faux")

	check.NoError(li(c))

	check.Equal(http.StatusTemporaryRedirect, rec.Result().StatusCode)

	cookies := rec.Result().Cookies()
	check.Len(cookies, 1)

	rawRedirectURL := rec.Header().Get("Location")
	check.NotEmpty(rawRedirectURL)

	qu, err := url.Parse(rawRedirectURL)
	check.NoError(err)

	q := qu.Query()

	// Use has logged in, is redirected to the callback page.

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/callback/faux?%s", q.Encode()), nil)
	req.AddCookie(cookies[0])

	rec = httptest.NewRecorder()

	sess := faux.Session{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}

	sn, err := store.Get(req, gothic.SessionName)
	check.NoError(err)
	sn.Values["faux"] = gzipString(sess.Marshal())
	check.NoError(sn.Save(req, rec))

	c = e.NewContext(req, rec)

	c.SetPath(fmt.Sprintf("/auth/callback/:provider?%s", q.Encode()))
	c.SetParamNames("provider")
	c.SetParamValues("faux")

	check.NoError(cb(c))

	check.Equal(http.StatusTemporaryRedirect, rec.Result().StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/whoami", nil)

	for _, c := range rec.Result().Cookies() {
		req.AddCookie(c)
	}

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	check.NoError(wai(c))

	check.Equal(http.StatusOK, rec.Result().StatusCode)

	data, err = io.ReadAll(rec.Result().Body)
	check.NoError(err)
	check.NotEmpty(data)

	t.Log(string(data))

	check.NoError(rec.Result().Body.Close())

	// Test Refetch

	req = httptest.NewRequest(http.MethodGet, "/auth/refetch/faux", nil)
	for _, c := range rec.Result().Cookies() {
		req.AddCookie(c)
	}
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	c.SetPath("/auth/refetch/:provider")
	c.SetParamNames("provider")
	c.SetParamValues("faux")

	check.NoError(refe(c))

	check.Equal(http.StatusTemporaryRedirect, rec.Result().StatusCode)

	qu, err = url.Parse(rec.Header().Get("Location"))
	check.NoError(err)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/callback/faux?%s", qu.Query().Encode()), nil)

	for _, c := range rec.Result().Cookies() {
		req.AddCookie(c)
	}
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	c.SetPath("/auth/callback/:provider")
	c.SetParamNames("provider")
	c.SetParamValues("faux")

	check.NoError(cb(c))

	check.Equal(http.StatusTemporaryRedirect, rec.Result().StatusCode)
	check.Equal(a.paths.afterLogin, rec.Header().Get("Location"))

}

func TestEr(t *testing.T) {

	check := require.New(t)

	port := 3215
	timeout := 10 * time.Second

	e := echo.New()
	sm := scs.New()
	sl := slogt.New(t)

	e.Use(session.LoadAndSave(sm))
	var err error

	gothic.Store = sessions.NewCookieStore([]byte("insecure_key"))

	check.NoError(New(e,
		SetDatabase(&dummy.DB{}),
		SetSessions(sm),
		SetLogger(sl),
		WithProvider(provider.OpenID, "kbyuFDidLLm280LIwVFiazOqjO3ty8KH", "60Op4HFM0I8ajz0WdiStAbziZ-VFQttXuxixHHs2R7r7-CW8GR79l-mmLqMhc-Sa", "https://openidconnect.net/callback", "https://samples.auth0.com/.well-known/openid-configuration"),
	))

	go func() {
		e.Start(fmt.Sprintf(":%d", port))
	}()

	time.Sleep(300 * time.Millisecond)

	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/auth/login/openid-connect", port), nil)
	check.NoError(err)

	resp, err := client.Do(req)
	check.NoError(err)

	check.NotEqual(http.StatusMethodNotAllowed, resp.StatusCode)

	// Shutdown the server

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	check.NoError(e.Shutdown(ctx))

}
