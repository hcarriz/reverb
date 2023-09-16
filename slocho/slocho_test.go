package slocho

import (
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/require"
)

func TestSlocho(t *testing.T) {

	check := require.New(t)

	l := slogt.New(t)

	port := 9081

	mw, err := New(Logger(l), WithHeaders("Authorization"), Level(slog.LevelError))
	check.NoError(err)

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(mw)
	e.GET("/", func(c echo.Context) error { return c.NoContent(http.StatusOK) })

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
			l.Error("unable to start server")
		}
	}()

	time.Sleep(300 * time.Millisecond)

	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/", port), nil)
	check.NoError(err)

	req.Header.Add("Authorization", "None")
	req.Header.Add("Authorization", "Bearer")

	resp, err := client.Do(req)
	check.NoError(err)
	check.Equal(http.StatusOK, resp.StatusCode)

}
