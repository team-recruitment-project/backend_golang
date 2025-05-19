package main

import (
	"backend_golang/internal/controller"
	"backend_golang/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	teamService := service.NewTeamService()
	teamController := controller.NewTeamController(teamService)

	app := gin.Default()

	app.POST("/v1/teams", teamController.MakeTeam)

	app.Run(":8080")
}
