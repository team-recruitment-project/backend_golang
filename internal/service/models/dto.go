package models

import "backend_golang/internal/models"

type CreateTeam struct {
	TeamName    string
	Description string
	Headcount   int8
	Vacancies   []models.Vacancy
}
