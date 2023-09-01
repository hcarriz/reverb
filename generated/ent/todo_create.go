// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/hcarriz/reverb/generated/ent/todo"
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TodoCreate is the builder for creating a Todo entity.
type TodoCreate struct {
	config
	mutation *TodoMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetContent sets the "content" field.
func (tc *TodoCreate) SetContent(s string) *TodoCreate {
	tc.mutation.SetContent(s)
	return tc
}

// SetDue sets the "due" field.
func (tc *TodoCreate) SetDue(t time.Time) *TodoCreate {
	tc.mutation.SetDue(t)
	return tc
}

// SetNillableDue sets the "due" field if the given value is not nil.
func (tc *TodoCreate) SetNillableDue(t *time.Time) *TodoCreate {
	if t != nil {
		tc.SetDue(*t)
	}
	return tc
}

// SetPriority sets the "priority" field.
func (tc *TodoCreate) SetPriority(t todo.Priority) *TodoCreate {
	tc.mutation.SetPriority(t)
	return tc
}

// SetNillablePriority sets the "priority" field if the given value is not nil.
func (tc *TodoCreate) SetNillablePriority(t *todo.Priority) *TodoCreate {
	if t != nil {
		tc.SetPriority(*t)
	}
	return tc
}

// Mutation returns the TodoMutation object of the builder.
func (tc *TodoCreate) Mutation() *TodoMutation {
	return tc.mutation
}

// Save creates the Todo in the database.
func (tc *TodoCreate) Save(ctx context.Context) (*Todo, error) {
	tc.defaults()
	return withHooks[*Todo, TodoMutation](ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TodoCreate) SaveX(ctx context.Context) *Todo {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TodoCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TodoCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TodoCreate) defaults() {
	if _, ok := tc.mutation.Priority(); !ok {
		v := todo.DefaultPriority
		tc.mutation.SetPriority(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TodoCreate) check() error {
	if _, ok := tc.mutation.Content(); !ok {
		return &ValidationError{Name: "content", err: errors.New(`ent: missing required field "Todo.content"`)}
	}
	if _, ok := tc.mutation.Priority(); !ok {
		return &ValidationError{Name: "priority", err: errors.New(`ent: missing required field "Todo.priority"`)}
	}
	if v, ok := tc.mutation.Priority(); ok {
		if err := todo.PriorityValidator(v); err != nil {
			return &ValidationError{Name: "priority", err: fmt.Errorf(`ent: validator failed for field "Todo.priority": %w`, err)}
		}
	}
	return nil
}

func (tc *TodoCreate) sqlSave(ctx context.Context) (*Todo, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TodoCreate) createSpec() (*Todo, *sqlgraph.CreateSpec) {
	var (
		_node = &Todo{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(todo.Table, sqlgraph.NewFieldSpec(todo.FieldID, field.TypeInt))
	)
	_spec.Schema = tc.schemaConfig.Todo
	_spec.OnConflict = tc.conflict
	if value, ok := tc.mutation.Content(); ok {
		_spec.SetField(todo.FieldContent, field.TypeString, value)
		_node.Content = value
	}
	if value, ok := tc.mutation.Due(); ok {
		_spec.SetField(todo.FieldDue, field.TypeTime, value)
		_node.Due = value
	}
	if value, ok := tc.mutation.Priority(); ok {
		_spec.SetField(todo.FieldPriority, field.TypeEnum, value)
		_node.Priority = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Todo.Create().
//		SetContent(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TodoUpsert) {
//			SetContent(v+v).
//		}).
//		Exec(ctx)
func (tc *TodoCreate) OnConflict(opts ...sql.ConflictOption) *TodoUpsertOne {
	tc.conflict = opts
	return &TodoUpsertOne{
		create: tc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Todo.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tc *TodoCreate) OnConflictColumns(columns ...string) *TodoUpsertOne {
	tc.conflict = append(tc.conflict, sql.ConflictColumns(columns...))
	return &TodoUpsertOne{
		create: tc,
	}
}

type (
	// TodoUpsertOne is the builder for "upsert"-ing
	//  one Todo node.
	TodoUpsertOne struct {
		create *TodoCreate
	}

	// TodoUpsert is the "OnConflict" setter.
	TodoUpsert struct {
		*sql.UpdateSet
	}
)

// SetContent sets the "content" field.
func (u *TodoUpsert) SetContent(v string) *TodoUpsert {
	u.Set(todo.FieldContent, v)
	return u
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *TodoUpsert) UpdateContent() *TodoUpsert {
	u.SetExcluded(todo.FieldContent)
	return u
}

// SetDue sets the "due" field.
func (u *TodoUpsert) SetDue(v time.Time) *TodoUpsert {
	u.Set(todo.FieldDue, v)
	return u
}

// UpdateDue sets the "due" field to the value that was provided on create.
func (u *TodoUpsert) UpdateDue() *TodoUpsert {
	u.SetExcluded(todo.FieldDue)
	return u
}

// ClearDue clears the value of the "due" field.
func (u *TodoUpsert) ClearDue() *TodoUpsert {
	u.SetNull(todo.FieldDue)
	return u
}

// SetPriority sets the "priority" field.
func (u *TodoUpsert) SetPriority(v todo.Priority) *TodoUpsert {
	u.Set(todo.FieldPriority, v)
	return u
}

// UpdatePriority sets the "priority" field to the value that was provided on create.
func (u *TodoUpsert) UpdatePriority() *TodoUpsert {
	u.SetExcluded(todo.FieldPriority)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Todo.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *TodoUpsertOne) UpdateNewValues() *TodoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Todo.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *TodoUpsertOne) Ignore() *TodoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TodoUpsertOne) DoNothing() *TodoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TodoCreate.OnConflict
// documentation for more info.
func (u *TodoUpsertOne) Update(set func(*TodoUpsert)) *TodoUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TodoUpsert{UpdateSet: update})
	}))
	return u
}

// SetContent sets the "content" field.
func (u *TodoUpsertOne) SetContent(v string) *TodoUpsertOne {
	return u.Update(func(s *TodoUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *TodoUpsertOne) UpdateContent() *TodoUpsertOne {
	return u.Update(func(s *TodoUpsert) {
		s.UpdateContent()
	})
}

// SetDue sets the "due" field.
func (u *TodoUpsertOne) SetDue(v time.Time) *TodoUpsertOne {
	return u.Update(func(s *TodoUpsert) {
		s.SetDue(v)
	})
}

// UpdateDue sets the "due" field to the value that was provided on create.
func (u *TodoUpsertOne) UpdateDue() *TodoUpsertOne {
	return u.Update(func(s *TodoUpsert) {
		s.UpdateDue()
	})
}

// ClearDue clears the value of the "due" field.
func (u *TodoUpsertOne) ClearDue() *TodoUpsertOne {
	return u.Update(func(s *TodoUpsert) {
		s.ClearDue()
	})
}

// SetPriority sets the "priority" field.
func (u *TodoUpsertOne) SetPriority(v todo.Priority) *TodoUpsertOne {
	return u.Update(func(s *TodoUpsert) {
		s.SetPriority(v)
	})
}

// UpdatePriority sets the "priority" field to the value that was provided on create.
func (u *TodoUpsertOne) UpdatePriority() *TodoUpsertOne {
	return u.Update(func(s *TodoUpsert) {
		s.UpdatePriority()
	})
}

// Exec executes the query.
func (u *TodoUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for TodoCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TodoUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *TodoUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *TodoUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// TodoCreateBulk is the builder for creating many Todo entities in bulk.
type TodoCreateBulk struct {
	config
	builders []*TodoCreate
	conflict []sql.ConflictOption
}

// Save creates the Todo entities in the database.
func (tcb *TodoCreateBulk) Save(ctx context.Context) ([]*Todo, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Todo, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TodoMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = tcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TodoCreateBulk) SaveX(ctx context.Context) []*Todo {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TodoCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TodoCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Todo.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TodoUpsert) {
//			SetContent(v+v).
//		}).
//		Exec(ctx)
func (tcb *TodoCreateBulk) OnConflict(opts ...sql.ConflictOption) *TodoUpsertBulk {
	tcb.conflict = opts
	return &TodoUpsertBulk{
		create: tcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Todo.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tcb *TodoCreateBulk) OnConflictColumns(columns ...string) *TodoUpsertBulk {
	tcb.conflict = append(tcb.conflict, sql.ConflictColumns(columns...))
	return &TodoUpsertBulk{
		create: tcb,
	}
}

// TodoUpsertBulk is the builder for "upsert"-ing
// a bulk of Todo nodes.
type TodoUpsertBulk struct {
	create *TodoCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Todo.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *TodoUpsertBulk) UpdateNewValues() *TodoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Todo.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *TodoUpsertBulk) Ignore() *TodoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TodoUpsertBulk) DoNothing() *TodoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TodoCreateBulk.OnConflict
// documentation for more info.
func (u *TodoUpsertBulk) Update(set func(*TodoUpsert)) *TodoUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TodoUpsert{UpdateSet: update})
	}))
	return u
}

// SetContent sets the "content" field.
func (u *TodoUpsertBulk) SetContent(v string) *TodoUpsertBulk {
	return u.Update(func(s *TodoUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *TodoUpsertBulk) UpdateContent() *TodoUpsertBulk {
	return u.Update(func(s *TodoUpsert) {
		s.UpdateContent()
	})
}

// SetDue sets the "due" field.
func (u *TodoUpsertBulk) SetDue(v time.Time) *TodoUpsertBulk {
	return u.Update(func(s *TodoUpsert) {
		s.SetDue(v)
	})
}

// UpdateDue sets the "due" field to the value that was provided on create.
func (u *TodoUpsertBulk) UpdateDue() *TodoUpsertBulk {
	return u.Update(func(s *TodoUpsert) {
		s.UpdateDue()
	})
}

// ClearDue clears the value of the "due" field.
func (u *TodoUpsertBulk) ClearDue() *TodoUpsertBulk {
	return u.Update(func(s *TodoUpsert) {
		s.ClearDue()
	})
}

// SetPriority sets the "priority" field.
func (u *TodoUpsertBulk) SetPriority(v todo.Priority) *TodoUpsertBulk {
	return u.Update(func(s *TodoUpsert) {
		s.SetPriority(v)
	})
}

// UpdatePriority sets the "priority" field to the value that was provided on create.
func (u *TodoUpsertBulk) UpdatePriority() *TodoUpsertBulk {
	return u.Update(func(s *TodoUpsert) {
		s.UpdatePriority()
	})
}

// Exec executes the query.
func (u *TodoUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the TodoCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for TodoCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TodoUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
