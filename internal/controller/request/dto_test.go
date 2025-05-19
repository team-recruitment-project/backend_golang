package request

import (
	"backend_golang/internal/models"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestMakeTeamRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     MakeTeamRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: MakeTeamRequest{
				TeamName:    "Test Team",
				Description: "Test Description",
				Headcount:   5,
				Vacancies: []models.Vacancy{
					models.NewVacancy(models.Backend, 2),
				},
			},
			wantErr: false,
		},
		{
			name: "empty team name",
			req: MakeTeamRequest{
				TeamName:    "",
				Description: "Test Description",
				Headcount:   5,
				Vacancies: []models.Vacancy{
					models.NewVacancy(models.Backend, 2),
				},
			},
			wantErr: true,
		},
		{
			name: "whitespace team name",
			req: MakeTeamRequest{
				TeamName:    "   ",
				Description: "Test Description",
				Headcount:   5,
				Vacancies: []models.Vacancy{
					models.NewVacancy(models.Backend, 2),
				},
			},
			wantErr: true,
		},
		{
			name: "empty description",
			req: MakeTeamRequest{
				TeamName:    "Test Team",
				Description: "",
				Headcount:   5,
				Vacancies: []models.Vacancy{
					models.NewVacancy(models.Backend, 2),
				},
			},
			wantErr: true,
		},
		{
			name: "whitespace description",
			req: MakeTeamRequest{
				TeamName:    "Test Team",
				Description: "   ",
				Headcount:   5,
				Vacancies: []models.Vacancy{
					models.NewVacancy(models.Backend, 2),
				},
			},
			wantErr: true,
		},
		{
			name: "empty vacancies",
			req: MakeTeamRequest{
				TeamName:    "Test Team",
				Description: "Test Description",
				Headcount:   5,
				Vacancies:   []models.Vacancy{},
			},
			wantErr: true,
		},
		{
			name: "nil vacancies",
			req: MakeTeamRequest{
				TeamName:    "Test Team",
				Description: "Test Description",
				Headcount:   5,
				Vacancies:   nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				// Check if the error is a validation error
				_, ok := err.(validator.ValidationErrors)
				assert.True(t, ok, "Error should be a ValidationErrors type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
