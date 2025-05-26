package models

import "backend_golang/internal/models"

type CreateTeam struct {
	MemberID    string
	TeamName    string
	Description string
	Headcount   int8
	Vacancies   []models.Vacancy
	Skills      []string
}

type RegisterAnnouncement struct {
	Title   string
	Content string
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
