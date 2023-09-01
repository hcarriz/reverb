// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/hcarriz/reverb/generated/ent/todo"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/entql"
	"entgo.io/ent/schema/field"
)

// schemaGraph holds a representation of ent/schema at runtime.
var schemaGraph = func() *sqlgraph.Schema {
	graph := &sqlgraph.Schema{Nodes: make([]*sqlgraph.Node, 1)}
	graph.Nodes[0] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   todo.Table,
			Columns: todo.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: todo.FieldID,
			},
		},
		Type: "Todo",
		Fields: map[string]*sqlgraph.FieldSpec{
			todo.FieldContent:  {Type: field.TypeString, Column: todo.FieldContent},
			todo.FieldDue:      {Type: field.TypeTime, Column: todo.FieldDue},
			todo.FieldPriority: {Type: field.TypeEnum, Column: todo.FieldPriority},
		},
	}
	return graph
}()

// predicateAdder wraps the addPredicate method.
// All update, update-one and query builders implement this interface.
type predicateAdder interface {
	addPredicate(func(s *sql.Selector))
}

// addPredicate implements the predicateAdder interface.
func (tq *TodoQuery) addPredicate(pred func(s *sql.Selector)) {
	tq.predicates = append(tq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the TodoQuery builder.
func (tq *TodoQuery) Filter() *TodoFilter {
	return &TodoFilter{config: tq.config, predicateAdder: tq}
}

// addPredicate implements the predicateAdder interface.
func (m *TodoMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the TodoMutation builder.
func (m *TodoMutation) Filter() *TodoFilter {
	return &TodoFilter{config: m.config, predicateAdder: m}
}

// TodoFilter provides a generic filtering capability at runtime for TodoQuery.
type TodoFilter struct {
	predicateAdder
	config
}

// Where applies the entql predicate on the query filter.
func (f *TodoFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[0].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *TodoFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(todo.FieldID))
}

// WhereContent applies the entql string predicate on the content field.
func (f *TodoFilter) WhereContent(p entql.StringP) {
	f.Where(p.Field(todo.FieldContent))
}

// WhereDue applies the entql time.Time predicate on the due field.
func (f *TodoFilter) WhereDue(p entql.TimeP) {
	f.Where(p.Field(todo.FieldDue))
}

// WherePriority applies the entql string predicate on the priority field.
func (f *TodoFilter) WherePriority(p entql.StringP) {
	f.Where(p.Field(todo.FieldPriority))
}
