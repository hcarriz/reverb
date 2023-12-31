// Code generated by ent, DO NOT EDIT.

package todo

import (
	"github.com/hcarriz/reverb/generated/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Todo {
	return predicate.Todo(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Todo {
	return predicate.Todo(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Todo {
	return predicate.Todo(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Todo {
	return predicate.Todo(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Todo {
	return predicate.Todo(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Todo {
	return predicate.Todo(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Todo {
	return predicate.Todo(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Todo {
	return predicate.Todo(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Todo {
	return predicate.Todo(sql.FieldLTE(FieldID, id))
}

// Content applies equality check predicate on the "content" field. It's identical to ContentEQ.
func Content(v string) predicate.Todo {
	return predicate.Todo(sql.FieldEQ(FieldContent, v))
}

// Due applies equality check predicate on the "due" field. It's identical to DueEQ.
func Due(v time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldEQ(FieldDue, v))
}

// ContentEQ applies the EQ predicate on the "content" field.
func ContentEQ(v string) predicate.Todo {
	return predicate.Todo(sql.FieldEQ(FieldContent, v))
}

// ContentNEQ applies the NEQ predicate on the "content" field.
func ContentNEQ(v string) predicate.Todo {
	return predicate.Todo(sql.FieldNEQ(FieldContent, v))
}

// ContentIn applies the In predicate on the "content" field.
func ContentIn(vs ...string) predicate.Todo {
	return predicate.Todo(sql.FieldIn(FieldContent, vs...))
}

// ContentNotIn applies the NotIn predicate on the "content" field.
func ContentNotIn(vs ...string) predicate.Todo {
	return predicate.Todo(sql.FieldNotIn(FieldContent, vs...))
}

// ContentGT applies the GT predicate on the "content" field.
func ContentGT(v string) predicate.Todo {
	return predicate.Todo(sql.FieldGT(FieldContent, v))
}

// ContentGTE applies the GTE predicate on the "content" field.
func ContentGTE(v string) predicate.Todo {
	return predicate.Todo(sql.FieldGTE(FieldContent, v))
}

// ContentLT applies the LT predicate on the "content" field.
func ContentLT(v string) predicate.Todo {
	return predicate.Todo(sql.FieldLT(FieldContent, v))
}

// ContentLTE applies the LTE predicate on the "content" field.
func ContentLTE(v string) predicate.Todo {
	return predicate.Todo(sql.FieldLTE(FieldContent, v))
}

// ContentContains applies the Contains predicate on the "content" field.
func ContentContains(v string) predicate.Todo {
	return predicate.Todo(sql.FieldContains(FieldContent, v))
}

// ContentHasPrefix applies the HasPrefix predicate on the "content" field.
func ContentHasPrefix(v string) predicate.Todo {
	return predicate.Todo(sql.FieldHasPrefix(FieldContent, v))
}

// ContentHasSuffix applies the HasSuffix predicate on the "content" field.
func ContentHasSuffix(v string) predicate.Todo {
	return predicate.Todo(sql.FieldHasSuffix(FieldContent, v))
}

// ContentEqualFold applies the EqualFold predicate on the "content" field.
func ContentEqualFold(v string) predicate.Todo {
	return predicate.Todo(sql.FieldEqualFold(FieldContent, v))
}

// ContentContainsFold applies the ContainsFold predicate on the "content" field.
func ContentContainsFold(v string) predicate.Todo {
	return predicate.Todo(sql.FieldContainsFold(FieldContent, v))
}

// DueEQ applies the EQ predicate on the "due" field.
func DueEQ(v time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldEQ(FieldDue, v))
}

// DueNEQ applies the NEQ predicate on the "due" field.
func DueNEQ(v time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldNEQ(FieldDue, v))
}

// DueIn applies the In predicate on the "due" field.
func DueIn(vs ...time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldIn(FieldDue, vs...))
}

// DueNotIn applies the NotIn predicate on the "due" field.
func DueNotIn(vs ...time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldNotIn(FieldDue, vs...))
}

// DueGT applies the GT predicate on the "due" field.
func DueGT(v time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldGT(FieldDue, v))
}

// DueGTE applies the GTE predicate on the "due" field.
func DueGTE(v time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldGTE(FieldDue, v))
}

// DueLT applies the LT predicate on the "due" field.
func DueLT(v time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldLT(FieldDue, v))
}

// DueLTE applies the LTE predicate on the "due" field.
func DueLTE(v time.Time) predicate.Todo {
	return predicate.Todo(sql.FieldLTE(FieldDue, v))
}

// DueIsNil applies the IsNil predicate on the "due" field.
func DueIsNil() predicate.Todo {
	return predicate.Todo(sql.FieldIsNull(FieldDue))
}

// DueNotNil applies the NotNil predicate on the "due" field.
func DueNotNil() predicate.Todo {
	return predicate.Todo(sql.FieldNotNull(FieldDue))
}

// PriorityEQ applies the EQ predicate on the "priority" field.
func PriorityEQ(v Priority) predicate.Todo {
	return predicate.Todo(sql.FieldEQ(FieldPriority, v))
}

// PriorityNEQ applies the NEQ predicate on the "priority" field.
func PriorityNEQ(v Priority) predicate.Todo {
	return predicate.Todo(sql.FieldNEQ(FieldPriority, v))
}

// PriorityIn applies the In predicate on the "priority" field.
func PriorityIn(vs ...Priority) predicate.Todo {
	return predicate.Todo(sql.FieldIn(FieldPriority, vs...))
}

// PriorityNotIn applies the NotIn predicate on the "priority" field.
func PriorityNotIn(vs ...Priority) predicate.Todo {
	return predicate.Todo(sql.FieldNotIn(FieldPriority, vs...))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Todo) predicate.Todo {
	return predicate.Todo(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Todo) predicate.Todo {
	return predicate.Todo(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Todo) predicate.Todo {
	return predicate.Todo(func(s *sql.Selector) {
		p(s.Not())
	})
}
