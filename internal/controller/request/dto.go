package request

import (
	"backend_golang/internal/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MakeTeamRequest struct {
	TeamName    string           `json:"teamName" validate:"required,min=1,notblank"`
	Description string           `json:"description" validate:"required,min=1,notblank"`
	Headcount   int8             `json:"headcount"`
	Vacancies   []models.Vacancy `json:"vacancies" validate:"required,min=1,dive"`
	Skills      []string         `json:"skills" validate:"required,min=1,dive,notblank"`
}

type PostAnnouncement struct {
	TeamID  int    `json:"teamID" validate:"required,min=1,notblank"`
	Title   string `json:"title" validate:"required,min=1,notblank"`
	Content string `json:"content" validate:"required,min=1,notblank"`
}

type SignUpRequest struct {
	Bio           string `json:"bio" validate:"required,min=1,notblank"`
	PreferredRole string `json:"preferredRole" validate:"required,min=1,notblank"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("notblank", validateNotBlank)
}

func validateNotBlank(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return strings.TrimSpace(value) != ""
}

func (r *MakeTeamRequest) Validate() error {
	return validate.Struct(r)
}

func (r *PostAnnouncement) Validate() error {
	return validate.Struct(r)
}

func (r *SignUpRequest) Validate() error {
	return validate.Struct(r)
}
