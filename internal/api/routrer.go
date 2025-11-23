package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"gorm.io/gorm"

	"go-rest-api/internal/api/handlers"
	"go-rest-api/internal/db/repository"
	"go-rest-api/internal/services"
)

func NewRouter(db *gorm.DB, logger *slog.Logger) *gin.Engine {
	router := gin.Default()

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPullRequestRepository(db)

	teamService := services.NewTeamService(db, teamRepo, userRepo)
	userService := services.NewUserService(db, userRepo, prRepo)
	prService := services.NewPullRequestService(db, prRepo, userRepo, teamRepo)

	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)
	prHandler := handlers.NewPullRequestHandler(prService)

	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())

	router.POST("/team/add", teamHandler.CreateTeam)
	router.GET("/team/get", teamHandler.GetTeam)

	router.POST("/users/setIsActive", userHandler.SetIsActive)
	router.GET("/users/getReview", userHandler.GetUserReviews)

	router.POST("/pullRequest/create", prHandler.CreatePR)
	router.POST("/pullRequest/merge", prHandler.MergePR)
	router.POST("/pullRequest/reassign", prHandler.ReassignReviewer)

	return router
}
