package service

import (
	"backend_golang/internal/service/models"
	"log"
)

type TeamService interface {
	Create(createTeam models.CreateTeam)
}

type teamService struct {
}

func NewTeamService() TeamService {
	return &teamService{}
}

func (t *teamService) Create(createTeam models.CreateTeam) {
	// Save to db
	log.Println("createTeam >>>", createTeam)
}
