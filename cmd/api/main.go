package main

import (
	"fmt"
	"log"

	_ "github.com/arturhk05/go-auth-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/arturhk05/go-auth-api/database"
	"github.com/arturhk05/go-auth-api/internal/handlers"
	"github.com/arturhk05/go-auth-api/internal/middlewares"
	"github.com/arturhk05/go-auth-api/internal/repositories"
	"github.com/arturhk05/go-auth-api/internal/services"
	"github.com/gin-gonic/gin"
)

// @title           Auth API
// @version         1.0
// @description     Authentication API built with Go, Gin, and PostgreSQL. Currently with only JWT authentication, but more features are planned for the future.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db.Db)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, cfg)
	userService := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	r.Use(middlewares.SecurityHeadersMiddleware())

	r.Use(middlewares.CORSMiddleware(cfg))

	public := r.Group("/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh-token", authHandler.RefreshToken)
	}

	protected := r.Group("/user")
	protected.Use(middlewares.AuthMiddleware(cfg))
	{
		protected.GET("/me", userHandler.GetProfile)
	}

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
