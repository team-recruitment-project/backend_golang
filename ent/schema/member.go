package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Member holds the schema definition for the Member entity.
type Member struct {
	ent.Schema
}

// Fields of the Member.
func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.String("member_id").Unique(),
		field.String("email").Unique(),
		field.String("picture"),
		field.String("nickname"),
		field.Text("bio"),
		field.String("preferred_role"),
	}
}

// Edges of the Member.
func (Member) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("skills", Skill.Type).
			Ref("users"),
	}
}
