package domain

import "time"

type Announcement struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Team      *Team
}
