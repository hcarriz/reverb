package authentication

import (
	"log/slog"
	"net/url"
)

func slerr(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func cloneURL(u *url.URL) *url.URL {

	if u == nil {
		return nil
	}

	u2 := new(url.URL)
	*u2 = *u

	if u.User != nil {
		u2.User = new(url.Userinfo)
		*u2.User = *u.User
	}

	return u2

}
