package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.Text("content"),
		field.Time("due").Optional(),
		field.Enum("priority").Values("low", "medium", "high").Default("medium"),
	}
}

// Annotations of the Todo
func (Todo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Mutations(),
	}
}
