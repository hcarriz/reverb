package sessions

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

// Option
type Option interface {
	apply(*Store) error
}

type option func(*Store) error

func (o option) apply(s *Store) error {
	return o(s)
}

func Timeout(timeout time.Duration) Option {
	return option(func(s *Store) error {

		s.timeout = timeout

		return nil
	})
}

// Cleanup
func Cleanup(cleanup time.Duration) Option {
	return option(func(s *Store) error {
		s.cleanup = cleanup
		return nil
	})
}

// Database
func Database(db Connection) Option {
	return option(func(s *Store) error {
		s.db = db
		return nil
	})
}

// Logger
func Logger(logger *slog.Logger) Option {
	return option(func(s *Store) error {
		s.logger = logger
		return nil
	})
}

type Connection interface {
	Find(context.Context, string) ([]byte, error)
	Delete(context.Context, string) error
	Add(context.Context, string, []byte, time.Time) error
	All(context.Context) (map[string][]byte, error)
	DeleteOld(context.Context, time.Time) error
}

type Store struct {
	cleanup time.Duration
	timeout time.Duration
	stop    chan bool
	logger  *slog.Logger
	db      Connection
}

func New(opts ...Option) (*Store, error) {

	s := &Store{
		logger:  slog.Default(),
		timeout: 1 * time.Minute,
		cleanup: 5 * time.Minute,
		db:      NewDB(),
	}

	var err error
	for _, opt := range opts {
		err = errors.Join(opt.apply(s))
	}

	if err != nil {
		return nil, err
	}

	if s.cleanup > 0 {
		go s.startCleanup()
	}

	return s, nil
}

func (s *Store) Find(token string) ([]byte, bool, error) {

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.timeout))
	defer cancel()

	result, err := s.db.Find(ctx, token)
	if err != nil {
		s.logger.ErrorContext(ctx, "unable to find token", slog.String("token", token), slog.String("error", err.Error()))
		return nil, false, err
	}

	s.logger.DebugContext(ctx, "found token", slog.String("token", token))

	return result, true, nil

}

func (s *Store) Delete(token string) error {

	now := time.Now()

	ctx, cancel := context.WithDeadline(context.Background(), now.Add(s.timeout))
	defer cancel()

	if err := s.db.Delete(ctx, token); err != nil {
		s.logger.ErrorContext(ctx, "unable to delete token", slog.String("token", token), slog.String("error", err.Error()))
		return err
	}

	s.logger.DebugContext(ctx, "token deleted", slog.String("token", token), slog.Time("when", now))

	return nil

}

func (s *Store) Commit(token string, data []byte, expires time.Time) error {

	now := time.Now()

	ctx, cancel := context.WithDeadline(context.Background(), now.Add(s.timeout))
	defer cancel()

	if err := s.db.Add(ctx, token, data, expires); err != nil {
		s.logger.ErrorContext(ctx, "unable to add token", slog.String("token", token), slog.String("data", string(data)), slog.String("error", err.Error()))
		return err
	}

	s.logger.DebugContext(ctx, "commited token", slog.String("token", token), slog.String("data", string(data)), slog.Time("expires", expires), slog.Time("added", now))

	return nil

}

func (s *Store) All() (map[string][]byte, error) {

	now := time.Now()

	ctx, cancel := context.WithDeadline(context.Background(), now.Add(s.timeout))
	defer cancel()

	result, err := s.db.All(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "unable to find tokens", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.DebugContext(ctx, "found tokens", slog.Int("amount", len(result)))

	return result, nil

}

func (s *Store) startCleanup() {

	s.stop = make(chan bool)

	if s.cleanup > 0 {

		ticker := time.NewTicker(s.cleanup)

		for {
			select {
			case <-ticker.C:

				s.logger.Debug("deleting sessions")

				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.timeout))
				defer cancel()

				if err := s.db.DeleteOld(ctx, time.Now()); err != nil {
					s.logger.ErrorContext(ctx, "unable to delete old tokens", slog.String("error", err.Error()))
				}

			case <-s.stop:
				s.logger.Info("stopping ticker")
				ticker.Stop()
				return
			}
		}

	}

}

func (s *Store) StopCleanup() {
	if s.stop != nil {
		s.stop <- true
		s.logger.Info("cleanup stopped")
	}
}
