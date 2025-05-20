package main

import (
	"backend_golang/ent"
	"backend_golang/internal/controller"
	"backend_golang/internal/repository"
	"backend_golang/internal/service"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	client, err := ent.Open("mysql", "team:team@tcp(localhost:3306)/team?parseTime=true")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	teamRepository := repository.NewTeamRepository(client)
	teamService := service.NewTeamService(teamRepository)
	teamController := controller.NewTeamController(teamService)

	app := gin.Default()

	app.POST("/v1/teams", teamController.MakeTeam)
	app.DELETE("/v1/teams/:teamID", teamController.DeleteTeam)

	app.Run(":8080")
}
