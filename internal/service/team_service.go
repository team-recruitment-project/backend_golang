package service

import (
	"backend_golang/internal/domain"
	imodels "backend_golang/internal/models"
	"backend_golang/internal/repository"
	"backend_golang/internal/service/models"
	"context"
)

type TeamService interface {
	Create(ctx context.Context, createTeam models.CreateTeam) (int, error)
	Delete(ctx context.Context, teamID int) error
	GetTeam(ctx context.Context, teamID int) (*models.TeamResponse, error)
}

type teamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) TeamService {
	return &teamService{teamRepository: teamRepository}
}

func (t *teamService) Create(ctx context.Context, createTeam models.CreateTeam) (int, error) {

	positions := make([]domain.Position, len(createTeam.Vacancies))
	for i, vacancy := range createTeam.Vacancies {
		positions[i] = domain.Position{
			Role:    imodels.Role(vacancy.Role),
			Vacancy: vacancy.Vacancy,
		}
	}

	skills := make([]domain.Skill, len(createTeam.Skills))
	for i, skill := range createTeam.Skills {
		skills[i] = domain.Skill{
			Name: skill,
		}
	}
	team, err := t.teamRepository.CreateTeam(ctx, &domain.Team{
		CreatedBy:   createTeam.MemberID,
		Name:        createTeam.TeamName,
		Description: createTeam.Description,
		Headcount:   createTeam.Headcount,
		Positions:   positions,
		Skills:      skills,
	})
	if err != nil {
		return 0, err
	}

	return team.ID, nil
}

func (t *teamService) Delete(ctx context.Context, teamID int) error {
	_, err := t.teamRepository.FindByID(ctx, teamID)
	if err != nil {
		return err
	}

	err = t.teamRepository.DeleteTeam(ctx, teamID)
	if err != nil {
		return err
	}
	return nil
}

func (t *teamService) GetTeam(ctx context.Context, teamID int) (*models.TeamResponse, error) {
	team, err := t.teamRepository.FindByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	vacancies := make([]imodels.Vacancy, len(team.Positions))
	for i, position := range team.Positions {
		vacancies[i] = imodels.Vacancy{
			Role:    position.Role,
			Vacancy: position.Vacancy,
		}
	}

	members := make([]domain.Member, len(team.Members))
	for i, member := range team.Members {
		members[i] = domain.Member{
			ID:            member.ID,
			Email:         member.Email,
			Nickname:      member.Nickname,
			Picture:       member.Picture,
			Bio:           member.Bio,
			PreferredRole: member.PreferredRole,
		}
	}
	return &models.TeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		Headcount:   team.Headcount,
		Members:     members,
		Vacancies:   vacancies,
		Skills:      team.Skills,
	}, nil
}
