package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Skill holds the schema definition for the Skill entity.
type Skill struct {
	ent.Schema
}

// Fields of the SKill.
func (Skill) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
	}
}

// Edges of the SKill.
func (Skill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", Member.Type),
		edge.To("teams", Team.Type),
	}
}
