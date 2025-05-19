package service

import (
	"backend_golang/internal/repository"
	"backend_golang/internal/service/models"
)

type TeamService interface {
	Create(createTeam models.CreateTeam) error
}

type teamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) TeamService {
	return &teamService{teamRepository: teamRepository}
}

func (t *teamService) Create(createTeam models.CreateTeam) error {

	err := t.teamRepository.CreateTeam(createTeam)
	if err != nil {
		return err
	}

	return nil
}
