package models

import (
	"backend_golang/internal/domain"
	"backend_golang/internal/models"
	"time"
)

type CreateTeam struct {
	MemberID    string
	TeamName    string
	Description string
	Headcount   int8
	Vacancies   []models.Vacancy
	Skills      []string
}

type RegisterAnnouncement struct {
	TeamID   int
	MemberID string
	Title    string
	Content  string
}

type SignupMember struct {
	Bio           string
	PreferredRole models.Role
}

type LoginResponse struct {
	URL string `json:"url"`
}

type UserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Nickname      string `json:"nickname"`
	Picture       string `json:"picture"`
	Bio           string `json:"bio"`
	PreferredRole string `json:"preferred_role"`
	Transient     bool   `json:"transient"`
}

type TeamResponse struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Headcount   int8             `json:"headcount"`
	CreatedBy   string           `json:"created_by"`
	Members     []domain.Member  `json:"members"`
	Vacancies   []models.Vacancy `json:"vacancies"`
	Skills      []domain.Skill   `json:"skills"`
}

type AnnouncementResponse struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Team      *TeamResponse `json:"team"`
}
