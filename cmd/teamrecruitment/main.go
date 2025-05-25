package main

import (
	"backend_golang/cmd/middleware"
	"backend_golang/ent"
	"backend_golang/internal/controller"
	"backend_golang/internal/repository"
	"backend_golang/internal/service"
	"context"
	"log"

	"github.com/gin-contrib/cors"
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

	app := gin.Default()

	// Middleware
	app.Use(cors.Default())

	// Team
	teamRepository := repository.NewTeamRepository(client)
	teamService := service.NewTeamService(teamRepository)
	teamController := controller.NewTeamController(teamService)
	app.POST("/v1/teams", teamController.MakeTeam)
	app.DELETE("/v1/teams/:teamID", teamController.DeleteTeam)

	// Announcement
	announcementRepository := repository.NewAnnouncementRepository(client)
	announcementService := service.NewAnnouncementService(announcementRepository)
	announcementController := controller.NewAnnouncementController(announcementService)
	app.POST("/v1/announcements", announcementController.Announce)

	// Auth
	authRepository := repository.NewAuthRepository(client)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	app.GET("/v1/auth/login", authController.Login)
	app.GET("/login/oauth2/code/google", authController.GoogleCallback)
	app.POST("/v1/auth/signup", authController.Signup, middleware.Authentication())
	app.GET("/v1/me", authController.GetMember, middleware.Authentication())

	app.Run(":8080")
}
