package middlewares

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/arturhk05/go-auth-api/config"
	apperrors "github.com/arturhk05/go-auth-api/internal/errors"
	"github.com/arturhk05/go-auth-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateAccessToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			if errors.Is(err, apperrors.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			} else if errors.Is(err, apperrors.ErrInvalidToken) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			} else {
				log.Printf("token validation error: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
