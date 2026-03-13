package main

import (
	"fmt"
	"log"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/database"
	"github.com/arturhk05/go-auth-api/internal/handlers"
	"github.com/arturhk05/go-auth-api/internal/repositories"
	"github.com/arturhk05/go-auth-api/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	userRepo := repositories.NewUserRepository(db.Db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	public := r.Group("/")
	{
		public.POST("/register", authHandler.Register)
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
