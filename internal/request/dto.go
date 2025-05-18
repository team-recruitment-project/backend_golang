package request

import "backend_golang/internal/models"

type MakeTeamRequest struct {
	TeamName    string
	Description string
	headcount   int8
	Vacancies   []models.Vacancy
}
