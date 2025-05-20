package repository

import (
	"backend_golang/ent/enttest"
	"backend_golang/internal/models"
	servicemodels "backend_golang/internal/service/models"
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestTeamRepository(t *testing.T) (*teamRepository, func()) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	repo := NewTeamRepository(client).(*teamRepository)
	return repo, func() {
		client.Close()
	}
}

func TestCreateTeam_Success(t *testing.T) {
	repo, cleanup := setupTestTeamRepository(t)
	defer cleanup()

	ctx := context.Background()
	createTeam := servicemodels.CreateTeam{
		TeamName:    "Test Team",
		Description: "Test Description",
		Headcount:   int8(5),
		Vacancies: []models.Vacancy{
			{
				Role:    "Developer",
				Vacancy: 2,
			},
			{
				Role:    "Designer",
				Vacancy: 1,
			},
		},
	}

	team, err := repo.CreateTeam(ctx, createTeam)
	require.NoError(t, err)
	assert.NotNil(t, team)
	assert.Equal(t, createTeam.TeamName, team.Name)
	assert.Equal(t, createTeam.Description, team.Description)
	assert.Equal(t, createTeam.Headcount, team.Headcount)

	// Verify positions were created correctly
	positions, err := repo.client.Position.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, positions, 2)
}

func TestCreateTeam_InvalidData(t *testing.T) {
	repo, cleanup := setupTestTeamRepository(t)
	defer cleanup()

	ctx := context.Background()
	createTeam := servicemodels.CreateTeam{
		TeamName:    "", // Invalid: empty team name
		Description: "Test Description",
		Headcount:   int8(5),
		Vacancies: []models.Vacancy{
			{
				Role:    "Developer",
				Vacancy: 2,
			},
		},
	}

	team, err := repo.CreateTeam(ctx, createTeam)
	assert.Error(t, err)
	assert.Nil(t, team)

	// Verify no positions were created
	positions, err := repo.client.Position.Query().All(ctx)
	require.NoError(t, err)
	assert.Empty(t, positions)
}

func TestCreateTeam_RollbackOnError(t *testing.T) {
	repo, cleanup := setupTestTeamRepository(t)
	defer cleanup()

	ctx := context.Background()
	createTeam := servicemodels.CreateTeam{
		TeamName:    "Test Team",
		Description: "Test Description",
		Headcount:   int8(5),
		Vacancies: []models.Vacancy{
			{
				Role:    "Developer",
				Vacancy: 2,
			},
		},
	}

	// First create a valid team
	team, err := repo.CreateTeam(ctx, createTeam)
	require.NoError(t, err)
	require.NotNil(t, team)

	// Try to create another team with a duplicate name
	// This should fail inside the transaction and rollback
	duplicateTeam := servicemodels.CreateTeam{
		TeamName:    "Test Team", // Same name as the first team
		Description: "Another Description",
		Headcount:   int8(3),
		Vacancies: []models.Vacancy{
			{
				Role:    "Designer",
				Vacancy: 1,
			},
		},
	}

	team2, err := repo.CreateTeam(ctx, duplicateTeam)
	assert.Error(t, err)
	assert.Nil(t, team2)

	// Verify only the first team exists
	teams, err := repo.client.Team.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, teams, 1)
	assert.Equal(t, createTeam.TeamName, teams[0].Name)

	// Verify only positions for the first team exist
	positions, err := repo.client.Position.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, positions, 1)
	assert.Equal(t, "Developer", positions[0].Role)
}

func TestCreateTeam_TransactionIsolation(t *testing.T) {
	repo, cleanup := setupTestTeamRepository(t)
	defer cleanup()

	ctx := context.Background()

	// First create a valid team
	createTeam1 := servicemodels.CreateTeam{
		TeamName:    "Team 1",
		Description: "Description 1",
		Headcount:   int8(3),
		Vacancies: []models.Vacancy{
			{
				Role:    "Developer",
				Vacancy: 2,
			},
		},
	}

	team1, err := repo.CreateTeam(ctx, createTeam1)
	require.NoError(t, err)
	assert.NotNil(t, team1)

	// Try to create a second team with invalid data
	createTeam2 := servicemodels.CreateTeam{
		TeamName:    "", // Invalid: empty team name
		Description: "Description 2",
		Headcount:   int8(2),
		Vacancies: []models.Vacancy{
			{
				Role:    "Designer",
				Vacancy: 1,
			},
		},
	}

	team2, err := repo.CreateTeam(ctx, createTeam2)
	assert.Error(t, err)
	assert.Nil(t, team2)

	// Verify first team still exists and is unchanged
	teams, err := repo.client.Team.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, teams, 1)
	assert.Equal(t, createTeam1.TeamName, teams[0].Name)

	// Verify positions for first team still exist
	positions, err := repo.client.Position.Query().All(ctx)
	require.NoError(t, err)
	assert.Len(t, positions, 1)
}
