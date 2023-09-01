package graceful

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/require"
)

func single(_ any, err error) error {
	return err
}

// testServer just return an OK status upon a request to "/"
func testServer(port int) Service {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return context.Background()
		},
	}

	return Service{
		Start: server.ListenAndServe,
		Close: func(context.Context) error {
			return server.Close()
		},
	}
}

func TestErrors(t *testing.T) {

	check := require.New(t)

	check.Error(single(Multiple([]Service{})))
	check.NoError(single(Single(testServer(8931))))
	check.NoError(single(Single(testServer(8932), Logger(slog.Default()))))
	check.Error(single(Single(testServer(8933), Timeout(-1))))
	check.NoError(single(Single(testServer(8934), Timeout(1))))

	quit, err := (Single(Service{
		Start: func() error { return errors.New("just for code coverage") },
		Close: func(context.Context) error { return errors.New("just for code coverage") },
	}))
	check.NoError(err)
	quit <- syscall.SIGINT

}

func TestChannel(t *testing.T) {

	check := require.New(t)

	port := 8888

	addr := fmt.Sprintf("http://localhost:%d", port)

	quit, err := Single(testServer(port))
	check.NoError(err)

	client := http.DefaultClient

	resp, err := client.Get(addr)
	check.NoError(err)
	check.NotNil(resp)
	check.Equal(http.StatusOK, resp.StatusCode)

	check.NoError(resp.Body.Close())

	quit <- syscall.SIGINT

	time.Sleep(300 * time.Millisecond)

	_, err = client.Get(addr)
	check.Error(err)

}

func TestSingle(t *testing.T) {

	check := require.New(t)

	port := 8889

	addr := fmt.Sprintf("http://localhost:%d", port)

	_, err := Single(testServer(port), Logger(slogt.New(t)))
	check.NoError(err)

	client := http.DefaultClient

	resp, err := client.Get(addr)
	check.NoError(err)
	check.NotNil(resp)
	check.Equal(http.StatusOK, resp.StatusCode)

	check.NoError(resp.Body.Close())

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	time.Sleep(300 * time.Millisecond)

	_, err = client.Get(addr)
	check.Error(err)

}
