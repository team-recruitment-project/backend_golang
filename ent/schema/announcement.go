package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Announcement holds the schema definition for the Announcement entity.
type Announcement struct {
	ent.Schema
}

// Fields of the Announcement.
func (Announcement) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("content").NotEmpty(),
	}
}

// Edges of the Announcement.
func (Announcement) Edges() []ent.Edge {
	return nil
}
