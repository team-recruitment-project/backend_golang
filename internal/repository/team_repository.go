package repository

import (
	"backend_golang/ent"
	"backend_golang/internal/service/models"
	"context"
	"log"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, createTeam models.CreateTeam) error
}

type teamRepository struct {
	client *ent.Client
}

func NewTeamRepository(client *ent.Client) TeamRepository {
	return &teamRepository{
		client: client,
	}
}

func (t *teamRepository) CreateTeam(ctx context.Context, createTeam models.CreateTeam) error {
	positions := []*ent.Position{}
	for _, vacancy := range createTeam.Vacancies {
		savedPosition, err := t.client.Position.Create().
			SetRole(string(vacancy.Role)).
			SetVacancy(vacancy.Vacancy).
			Save(ctx)
		if err != nil {
			log.Println("error", err)
			return err
		}
		positions = append(positions, savedPosition)
	}
	log.Println("positions", positions)

	t.client.Team.Create().
		SetName(createTeam.TeamName).
		SetDescription(createTeam.Description).
		SetHeadcount(createTeam.Headcount).
		AddPositions(positions...).
		Save(ctx)
	return nil
}
