package repository

import (
	"backend_golang/ent"
	"backend_golang/ent/position"
	"backend_golang/ent/team"
	"backend_golang/internal/domain"
	"backend_golang/internal/service/models"
	"context"
	"fmt"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, createTeam models.CreateTeam) (*domain.Team, error)
	DeleteTeam(ctx context.Context, teamID int) error
	FindByID(ctx context.Context, teamID int) error
}

type teamRepository struct {
	client *ent.Client
	tx     *TransactionManager
}

func NewTeamRepository(client *ent.Client) TeamRepository {
	return &teamRepository{
		client: client,
		tx:     NewTransactionManager(client),
	}
}

func (t *teamRepository) CreateTeam(ctx context.Context, createTeam models.CreateTeam) (*domain.Team, error) {
	// TODO : createTeam 을 서비스의 dto 가 아니라 리포지토리단의 domain 모델로 변경
	// Validate input
	if createTeam.TeamName == "" {
		return nil, fmt.Errorf("team name cannot be empty")
	}
	for _, vacancy := range createTeam.Vacancies {
		if vacancy.Role == "" {
			return nil, fmt.Errorf("role cannot be empty")
		}
	}

	var result *domain.Team
	err := t.tx.WithTx(ctx, func(tx *ent.Tx) error {
		positions := []*ent.Position{}
		for _, vacancy := range createTeam.Vacancies {
			savedPosition, err := tx.Position.Create().
				SetRole(string(vacancy.Role)).
				SetVacancy(vacancy.Vacancy).
				Save(ctx)
			if err != nil {
				return err
			}
			positions = append(positions, savedPosition)
		}

		team, err := tx.Team.Create().
			SetName(createTeam.TeamName).
			SetDescription(createTeam.Description).
			SetHeadcount(createTeam.Headcount).
			AddPositions(positions...).
			Save(ctx)
		if err != nil {
			return err
		}

		positionIDs := make([]int, len(team.Edges.Positions))
		for i, position := range team.Edges.Positions {
			positionIDs[i] = position.ID
		}

		result = &domain.Team{
			ID:          team.ID,
			Name:        team.Name,
			Description: team.Description,
			Headcount:   team.Headcount,
			Positions:   positionIDs,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *teamRepository) DeleteTeam(ctx context.Context, teamID int) error {
	return t.tx.WithTx(ctx, func(tx *ent.Tx) error {
		// First, delete all positions associated with the team
		_, err := tx.Position.Delete().Where(
			position.HasTeamWith(team.ID(teamID)),
		).Exec(ctx)
		if err != nil {
			return err
		}

		// Then delete the team
		err = tx.Team.DeleteOneID(teamID).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

func (t *teamRepository) FindByID(ctx context.Context, teamID int) error {
	_, err := t.client.Team.Query().Where(team.ID(teamID)).First(ctx)
	if err != nil {
		return err
	}
	return nil
}
