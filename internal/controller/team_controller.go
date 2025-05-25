package controller

import (
	"backend_golang/internal/controller/request"
	"backend_golang/internal/models"
	"backend_golang/internal/service"
	smodels "backend_golang/internal/service/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TeamController interface {
	MakeTeam(c *gin.Context)
	DeleteTeam(c *gin.Context)
}

type teamController struct {
	teamService service.TeamService
}

func NewTeamController(teamService service.TeamService) TeamController {
	return &teamController{teamService: teamService}
}

func (t *teamController) MakeTeam(c *gin.Context) {
	req := &request.MakeTeamRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	teamID, err := t.teamService.Create(c,
		smodels.CreateTeam{
			TeamName:    req.TeamName,
			Description: req.Description,
			Headcount:   req.Headcount,
			Vacancies:   req.Vacancies,
			Skills:      req.Skills,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"teamID": teamID})
}

func (t *teamController) DeleteTeam(c *gin.Context) {
	log.Println("test >>", c.Param("teamID"))
	teamID, err := strconv.Atoi(c.Param("teamID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t.teamService.Delete(c, teamID)
}
