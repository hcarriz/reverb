package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"
)

var (
	ErrInvalidDuration = errors.New("invalid duration")
	ErrExistingAddress = errors.New("the address has already been set")
	ErrInvalidAddr     = errors.New("invalid address")
	ErrMissingServices = errors.New("missing services")
)

type config struct {
	timeout time.Duration
	logger  Log
}

type Log interface {
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

type Option interface {
	apply(*config) error
}

type option func(*config) error

func (o option) apply(c *config) error {
	return o(c)
}

func Timeout(duration time.Duration) Option {
	return option(func(c *config) error {

		if duration < 0 {
			return ErrInvalidDuration
		}

		c.timeout = duration

		return nil
	})
}

func Logger(logger Log) Option {
	return option(func(c *config) error {
		c.logger = logger
		return nil
	})
}

type Service struct {
	Start func() error
	Close func(context.Context) error
}

// Single allows for a single Service to operate.
func Single(service Service, opts ...Option) (chan os.Signal, error) {
	return Multiple([]Service{service}, opts...)
}

// Multiple allows for multiple services to operate.
func Multiple(services []Service, opts ...Option) (chan os.Signal, error) {

	cfg := config{
		timeout: 5 * time.Second,
		logger:  slog.Default(),
	}

	if len(services) < 1 {
		return nil, ErrMissingServices
	}

	var err error

	for _, opt := range opts {
		err = errors.Join(opt.apply(&cfg))
	}

	if err != nil {
		return nil, err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for _, service := range services {

		service := service

		go func() {

			<-quit

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.timeout))
			defer cancel()

			if err := service.Close(ctx); err != nil {
				cfg.logger.LogAttrs(ctx, slog.LevelError, "unable to close service", slog.String("error", err.Error()))
			}

		}()

		go func() {

			if err := service.Start(); err != nil {
				cfg.logger.LogAttrs(context.Background(), slog.LevelError, "service error", slog.String("error", err.Error()))
			}

		}()

	}

	return quit, nil

}
