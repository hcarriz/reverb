package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"
	"reverb/generated/ent"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input ent.CreateTodoInput) (*ent.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, input ent.UpdateTodoInput) (*ent.Todo, error) {
	panic(fmt.Errorf("not implemented: UpdateTodo - updateTodo"))
}

// DeleteTodo is the resolver for the deleteTodo field.
func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (*ent.Todo, error) {
	panic(fmt.Errorf("not implemented: DeleteTodo - deleteTodo"))
}

// Todo is the resolver for the todo field.
func (r *queryResolver) Todo(ctx context.Context, id *string) (*ent.Todo, error) {
	panic(fmt.Errorf("not implemented: Todo - todo"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
