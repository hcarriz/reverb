// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/hcarriz/reverb/generated/ent/predicate"
	"github.com/hcarriz/reverb/generated/ent/todo"
	"errors"
	"fmt"
	"time"
)

// TodoWhereInput represents a where input for filtering Todo queries.
type TodoWhereInput struct {
	Predicates []predicate.Todo  `json:"-"`
	Not        *TodoWhereInput   `json:"not,omitempty"`
	Or         []*TodoWhereInput `json:"or,omitempty"`
	And        []*TodoWhereInput `json:"and,omitempty"`

	// "id" field predicates.
	ID      *int  `json:"id,omitempty"`
	IDNEQ   *int  `json:"idNEQ,omitempty"`
	IDIn    []int `json:"idIn,omitempty"`
	IDNotIn []int `json:"idNotIn,omitempty"`
	IDGT    *int  `json:"idGT,omitempty"`
	IDGTE   *int  `json:"idGTE,omitempty"`
	IDLT    *int  `json:"idLT,omitempty"`
	IDLTE   *int  `json:"idLTE,omitempty"`

	// "content" field predicates.
	Content             *string  `json:"content,omitempty"`
	ContentNEQ          *string  `json:"contentNEQ,omitempty"`
	ContentIn           []string `json:"contentIn,omitempty"`
	ContentNotIn        []string `json:"contentNotIn,omitempty"`
	ContentGT           *string  `json:"contentGT,omitempty"`
	ContentGTE          *string  `json:"contentGTE,omitempty"`
	ContentLT           *string  `json:"contentLT,omitempty"`
	ContentLTE          *string  `json:"contentLTE,omitempty"`
	ContentContains     *string  `json:"contentContains,omitempty"`
	ContentHasPrefix    *string  `json:"contentHasPrefix,omitempty"`
	ContentHasSuffix    *string  `json:"contentHasSuffix,omitempty"`
	ContentEqualFold    *string  `json:"contentEqualFold,omitempty"`
	ContentContainsFold *string  `json:"contentContainsFold,omitempty"`

	// "due" field predicates.
	Due       *time.Time  `json:"due,omitempty"`
	DueNEQ    *time.Time  `json:"dueNEQ,omitempty"`
	DueIn     []time.Time `json:"dueIn,omitempty"`
	DueNotIn  []time.Time `json:"dueNotIn,omitempty"`
	DueGT     *time.Time  `json:"dueGT,omitempty"`
	DueGTE    *time.Time  `json:"dueGTE,omitempty"`
	DueLT     *time.Time  `json:"dueLT,omitempty"`
	DueLTE    *time.Time  `json:"dueLTE,omitempty"`
	DueIsNil  bool        `json:"dueIsNil,omitempty"`
	DueNotNil bool        `json:"dueNotNil,omitempty"`

	// "priority" field predicates.
	Priority      *todo.Priority  `json:"priority,omitempty"`
	PriorityNEQ   *todo.Priority  `json:"priorityNEQ,omitempty"`
	PriorityIn    []todo.Priority `json:"priorityIn,omitempty"`
	PriorityNotIn []todo.Priority `json:"priorityNotIn,omitempty"`
}

// AddPredicates adds custom predicates to the where input to be used during the filtering phase.
func (i *TodoWhereInput) AddPredicates(predicates ...predicate.Todo) {
	i.Predicates = append(i.Predicates, predicates...)
}

// Filter applies the TodoWhereInput filter on the TodoQuery builder.
func (i *TodoWhereInput) Filter(q *TodoQuery) (*TodoQuery, error) {
	if i == nil {
		return q, nil
	}
	p, err := i.P()
	if err != nil {
		if err == ErrEmptyTodoWhereInput {
			return q, nil
		}
		return nil, err
	}
	return q.Where(p), nil
}

// ErrEmptyTodoWhereInput is returned in case the TodoWhereInput is empty.
var ErrEmptyTodoWhereInput = errors.New("ent: empty predicate TodoWhereInput")

// P returns a predicate for filtering todos.
// An error is returned if the input is empty or invalid.
func (i *TodoWhereInput) P() (predicate.Todo, error) {
	var predicates []predicate.Todo
	if i.Not != nil {
		p, err := i.Not.P()
		if err != nil {
			return nil, fmt.Errorf("%w: field 'not'", err)
		}
		predicates = append(predicates, todo.Not(p))
	}
	switch n := len(i.Or); {
	case n == 1:
		p, err := i.Or[0].P()
		if err != nil {
			return nil, fmt.Errorf("%w: field 'or'", err)
		}
		predicates = append(predicates, p)
	case n > 1:
		or := make([]predicate.Todo, 0, n)
		for _, w := range i.Or {
			p, err := w.P()
			if err != nil {
				return nil, fmt.Errorf("%w: field 'or'", err)
			}
			or = append(or, p)
		}
		predicates = append(predicates, todo.Or(or...))
	}
	switch n := len(i.And); {
	case n == 1:
		p, err := i.And[0].P()
		if err != nil {
			return nil, fmt.Errorf("%w: field 'and'", err)
		}
		predicates = append(predicates, p)
	case n > 1:
		and := make([]predicate.Todo, 0, n)
		for _, w := range i.And {
			p, err := w.P()
			if err != nil {
				return nil, fmt.Errorf("%w: field 'and'", err)
			}
			and = append(and, p)
		}
		predicates = append(predicates, todo.And(and...))
	}
	predicates = append(predicates, i.Predicates...)
	if i.ID != nil {
		predicates = append(predicates, todo.IDEQ(*i.ID))
	}
	if i.IDNEQ != nil {
		predicates = append(predicates, todo.IDNEQ(*i.IDNEQ))
	}
	if len(i.IDIn) > 0 {
		predicates = append(predicates, todo.IDIn(i.IDIn...))
	}
	if len(i.IDNotIn) > 0 {
		predicates = append(predicates, todo.IDNotIn(i.IDNotIn...))
	}
	if i.IDGT != nil {
		predicates = append(predicates, todo.IDGT(*i.IDGT))
	}
	if i.IDGTE != nil {
		predicates = append(predicates, todo.IDGTE(*i.IDGTE))
	}
	if i.IDLT != nil {
		predicates = append(predicates, todo.IDLT(*i.IDLT))
	}
	if i.IDLTE != nil {
		predicates = append(predicates, todo.IDLTE(*i.IDLTE))
	}
	if i.Content != nil {
		predicates = append(predicates, todo.ContentEQ(*i.Content))
	}
	if i.ContentNEQ != nil {
		predicates = append(predicates, todo.ContentNEQ(*i.ContentNEQ))
	}
	if len(i.ContentIn) > 0 {
		predicates = append(predicates, todo.ContentIn(i.ContentIn...))
	}
	if len(i.ContentNotIn) > 0 {
		predicates = append(predicates, todo.ContentNotIn(i.ContentNotIn...))
	}
	if i.ContentGT != nil {
		predicates = append(predicates, todo.ContentGT(*i.ContentGT))
	}
	if i.ContentGTE != nil {
		predicates = append(predicates, todo.ContentGTE(*i.ContentGTE))
	}
	if i.ContentLT != nil {
		predicates = append(predicates, todo.ContentLT(*i.ContentLT))
	}
	if i.ContentLTE != nil {
		predicates = append(predicates, todo.ContentLTE(*i.ContentLTE))
	}
	if i.ContentContains != nil {
		predicates = append(predicates, todo.ContentContains(*i.ContentContains))
	}
	if i.ContentHasPrefix != nil {
		predicates = append(predicates, todo.ContentHasPrefix(*i.ContentHasPrefix))
	}
	if i.ContentHasSuffix != nil {
		predicates = append(predicates, todo.ContentHasSuffix(*i.ContentHasSuffix))
	}
	if i.ContentEqualFold != nil {
		predicates = append(predicates, todo.ContentEqualFold(*i.ContentEqualFold))
	}
	if i.ContentContainsFold != nil {
		predicates = append(predicates, todo.ContentContainsFold(*i.ContentContainsFold))
	}
	if i.Due != nil {
		predicates = append(predicates, todo.DueEQ(*i.Due))
	}
	if i.DueNEQ != nil {
		predicates = append(predicates, todo.DueNEQ(*i.DueNEQ))
	}
	if len(i.DueIn) > 0 {
		predicates = append(predicates, todo.DueIn(i.DueIn...))
	}
	if len(i.DueNotIn) > 0 {
		predicates = append(predicates, todo.DueNotIn(i.DueNotIn...))
	}
	if i.DueGT != nil {
		predicates = append(predicates, todo.DueGT(*i.DueGT))
	}
	if i.DueGTE != nil {
		predicates = append(predicates, todo.DueGTE(*i.DueGTE))
	}
	if i.DueLT != nil {
		predicates = append(predicates, todo.DueLT(*i.DueLT))
	}
	if i.DueLTE != nil {
		predicates = append(predicates, todo.DueLTE(*i.DueLTE))
	}
	if i.DueIsNil {
		predicates = append(predicates, todo.DueIsNil())
	}
	if i.DueNotNil {
		predicates = append(predicates, todo.DueNotNil())
	}
	if i.Priority != nil {
		predicates = append(predicates, todo.PriorityEQ(*i.Priority))
	}
	if i.PriorityNEQ != nil {
		predicates = append(predicates, todo.PriorityNEQ(*i.PriorityNEQ))
	}
	if len(i.PriorityIn) > 0 {
		predicates = append(predicates, todo.PriorityIn(i.PriorityIn...))
	}
	if len(i.PriorityNotIn) > 0 {
		predicates = append(predicates, todo.PriorityNotIn(i.PriorityNotIn...))
	}

	switch len(predicates) {
	case 0:
		return nil, ErrEmptyTodoWhereInput
	case 1:
		return predicates[0], nil
	default:
		return todo.And(predicates...), nil
	}
}
