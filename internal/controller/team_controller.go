package controller

import (
	"backend_golang/internal/models"
	"backend_golang/internal/request"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TeamController interface {
	MakeTeam(c *gin.Context)
}

type teamController struct {
}

func NewTeamController() TeamController {
	return &teamController{}
}

func (t *teamController) MakeTeam(c *gin.Context) {
	req := &request.MakeTeamRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		validationErrors := make([]models.ValidationError, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, models.NewValidationError(err))
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationErrors,
		})
		return
	}

	log.Println(req)
}
