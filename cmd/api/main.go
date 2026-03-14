package main

import (
	"fmt"
	"log"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/database"
	"github.com/arturhk05/go-auth-api/internal/handlers"
	"github.com/arturhk05/go-auth-api/internal/middlewares"
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
	defer db.Db.Close()

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	userRepo := repositories.NewUserRepository(db.Db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	r.Use(middlewares.SecurityHeadersMiddleware())

	r.Use(middlewares.CORSMiddleware(cfg))

	public := r.Group("/")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	if cfg.Server.TLSCertFile != "" && cfg.Server.TLSKeyFile != "" {
		log.Println("Starting server with TLS/HTTPS")
		if err := r.RunTLS(addr, cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		log.Println("Starting server with HTTP (no TLS configured)")
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
