package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Member holds the schema definition for the Member entity.
type Member struct {
	ent.Schema
}

// Fields of the Member.
func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("member_id"),
		field.String("email"),
		field.String("picture"),
		field.String("nickname"),
		field.Text("bio"),
		field.String("preferred_role"),
	}
}

// Edges of the Member.
func (Member) Edges() []ent.Edge {
	return nil
}
