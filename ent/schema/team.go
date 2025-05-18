package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Team holds the schema definition for the Team entity.
type Team struct {
	ent.Schema
}

// Fields of the Team.
func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("team_id"),
		field.String("name"),
		field.Text("description"),
		field.Int8("headcount"),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return nil
}
