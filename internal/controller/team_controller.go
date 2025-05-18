package controller

import (
	"backend_golang/internal/request"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, err)
		return
	}

	log.Println(req)
}
