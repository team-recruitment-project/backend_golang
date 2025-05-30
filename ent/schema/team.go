package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Team holds the schema definition for the Team entity.
type Team struct {
	ent.Schema
}

// Fields of the Team.
func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.Text("description"),
		field.Int8("headcount"),
		field.String("created_by").Unique().NotEmpty(),
	}
}

// Edges of the Team.
func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("positions", Position.Type),
		edge.To("members", Member.Type),
		edge.To("announcements", Announcement.Type),
		edge.From("skills", Skill.Type).
			Ref("teams"),
	}
}
