package grapher

import (
	"errors"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

type grapher struct {
	handler *handler.Server
}

type Option interface {
	apply(*grapher) error
}

type option func(*grapher) error

func (o option) apply(g *grapher) error {
	return o(g)
}

// Extensions
func WithExtensions(extensions ...graphql.HandlerExtension) Option {
	return option(func(g *grapher) error {

		for _, ext := range extensions {
			g.handler.Use(ext)
		}

		return nil
	})
}

func Ent(client entgql.TxOpener) Option {
	return WithExtensions(entgql.Transactioner{TxOpener: client})
}

func Cache(cache graphql.Cache) Option {
	return WithExtensions(extension.AutomaticPersistedQuery{Cache: cache})
}

// Transport
func WithTransports(transports ...graphql.Transport) Option {
	return option(func(g *grapher) error {

		for _, transport := range transports {
			g.handler.AddTransport(transport)
		}

		return nil
	})
}

// Default Transportation
func DefaultWebsockets() Option {
	return WithTransports(&transport.Websocket{})
}

func POST(in *transport.POST) Option {
	return WithTransports(in)
}

func DefaultPOST() Option {
	return POST(&transport.POST{})
}

func GET(in *transport.GET) Option {
	return WithTransports(in)
}

func DefaultGET() Option {
	return GET(&transport.GET{})
}

func MultipartForm(in *transport.MultipartForm) Option {
	return WithTransports(in)
}

func DefaultMultipartForm() Option {
	return MultipartForm(&transport.MultipartForm{})
}

func New(in *handler.Server, opts ...Option) error {

	g := &grapher{in}

	var err error

	for _, opt := range opts {
		err = errors.Join(opt.apply(g))
	}

	return err
}
