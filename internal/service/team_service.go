package service

import (
	"backend_golang/internal/repository"
	"backend_golang/internal/service/models"
	"context"
)

type TeamService interface {
	Create(ctx context.Context, createTeam models.CreateTeam) (int, error)
	Delete(ctx context.Context, teamID int) error
}

type teamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) TeamService {
	return &teamService{teamRepository: teamRepository}
}

func (t *teamService) Create(ctx context.Context, createTeam models.CreateTeam) (int, error) {
	team, err := t.teamRepository.CreateTeam(ctx, createTeam)
	if err != nil {
		return 0, err
	}

	return team.ID, nil
}

func (t *teamService) Delete(ctx context.Context, teamID int) error {
	err := t.teamRepository.FindByID(ctx, teamID)
	if err != nil {
		return err
	}

	err = t.teamRepository.DeleteTeam(ctx, teamID)
	if err != nil {
		return err
	}
	return nil
}
