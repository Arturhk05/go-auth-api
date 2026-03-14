package middlewares

import (
	"log"
	"time"

	"github.com/arturhk05/go-auth-api/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"HEAD",
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if cfg.Server.Environment == "production" {
		frontendURL := cfg.Server.FrontendURL
		if frontendURL != "" {
			corsConfig.AllowOrigins = []string{frontendURL}
		} else {
			log.Println("WARNING: GIN_MODE=release but FRONTEND_URL is not set. Set FRONTEND_URL to your frontend domain.")
		}
	} else {
		corsConfig.AllowOrigins = []string{
			"http://localhost:3000", // React dev server
			"http://localhost:5173", // Vite dev server
			"http://localhost:5174", // Vite dev server
			"http://localhost:8080", // API server itself
		}
	}

	return cors.New(corsConfig)
}
