package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// TransientMember holds the schema definition for the TransientMember entity.
type TransientMember struct {
	ent.Schema
}

// Fields of the TransientMember.
func (TransientMember) Fields() []ent.Field {
	return []ent.Field{
		field.String("transient_member_id").Unique(),
		field.String("email").Unique(),
		field.String("picture"),
		field.String("nickname"),
	}
}

// Edges of the TransientMember.
func (TransientMember) Edges() []ent.Edge {
	return nil
}
