package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Position holds the schema definition for the Position entity.
type Position struct {
	ent.Schema
}

// Fields of the Position.
func (Position) Fields() []ent.Field {
	return []ent.Field{
		// field.Int64("position_id"),
		field.Int("team_id").Optional(),
		field.String("role"),
		field.Int8("vacancy"),
	}
}

// Edges of the Position.
func (Position) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("team", Team.Type).
			Ref("positions").
			Unique().
			Field("team_id"),
	}
}
