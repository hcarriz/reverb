package viewer

import (
	"context"
)

type Value struct {
	v string
}

func (v Value) String() string {
	return v.v
}

var (
	ContextUserID = Value{"viewer_user_id"}
	ContextSystem = Value{"viewer_system"}
	ContextIP     = Value{"viewer_ip"}
)

type ID interface {
	~int | ~string
}

func getter[T any](ctx context.Context, thing Value) (T, bool) {
	if id, ok := ctx.Value(thing).(T); ok {
		return id, true
	}

	var empty T

	return empty, false
}

func setter[T any](ctx context.Context, thing Value, data T) context.Context {
	return context.WithValue(ctx, thing, data)
}

// Get
func Get[T ID](ctx context.Context, thing Value) (T, bool) {
	return getter[T](ctx, thing)
}

func Set[T ID](ctx context.Context, thing Value, data T) context.Context {
	return setter[T](ctx, thing, data)
}

// User ID
func GetUserID[T ID](ctx context.Context) (T, bool) {
	return Get[T](ctx, ContextUserID)
}

func SetUserID[T ID](ctx context.Context, id T) context.Context {
	return Set(ctx, ContextUserID, id)
}

// System
func SetSystem(ctx context.Context) context.Context {
	return setter(ctx, ContextSystem, true)
}

func IsSystem(ctx context.Context) bool {
	result, ok := getter[bool](ctx, ContextSystem)
	return ok && result
}

// IP Address

func SetAddress(ctx context.Context, ip string) context.Context {
	return Set(ctx, ContextIP, ip)
}

func GetAddress(ctx context.Context) string {

	if ip, ok := Get[string](ctx, ContextIP); ok {
		return ip
	}

	return "127.0.0.1"

}
