package handlers

import (
	"errors"
	"log"
	"net/http"

	apperrors "github.com/arturhk05/go-auth-api/internal/errors"
	"github.com/arturhk05/go-auth-api/internal/models"
	"github.com/arturhk05/go-auth-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      User Registration
// @Description  Register a new user and return access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "Registration credentials"
// @Success      201   {object}  models.AuthResponse
// @Failure      400   {object}  models.ErrorResponse  "Invalid request or validation failed"
// @Failure      409   {object}  models.ErrorResponse  "User already exists"
// @Failure      500   {object}  models.ErrorResponse  "Internal server error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "message": err.Error()})
		return
	}

	resp, err := h.authService.Register(req.Password, req.Email, req.Username)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		log.Printf("registration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and return access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Login credentials"
// @Success      200   {object}  models.AuthResponse
// @Failure      400   {object}  models.ErrorResponse  "Invalid request or validation failed"
// @Failure      401   {object}  models.ErrorResponse  "Invalid email or password"
// @Failure      403   {object}  models.ErrorResponse  "Account is inactive or locked"
// @Failure      500   {object}  models.ErrorResponse  "Internal server error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	resp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		if errors.Is(err, apperrors.ErrAccountInactive) {
			c.JSON(http.StatusForbidden, gin.H{"error": "account is inactive"})
			return
		}
		if errors.Is(err, apperrors.ErrAccountLocked) {
			c.JSON(http.StatusForbidden, gin.H{"error": "account is locked"})
			return
		}
		log.Printf("login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary      Refresh Auth Tokens
// @Description  Refresh access tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RefreshRequest  true  "Refresh token"
// @Success      200   {object}  models.AuthResponse
// @Failure      400   {object}  models.ErrorResponse  "Invalid request"
// @Failure      401   {object}  models.ErrorResponse  "Token expired, invalid, or revoked"
// @Failure      403   {object}  models.ErrorResponse  "Account is inactive or locked"
// @Failure      500   {object}  models.ErrorResponse  "Internal server error"
// @Router       /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}
		if errors.Is(err, apperrors.ErrInvalidToken) || errors.Is(err, apperrors.ErrTokenRevoked) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		if errors.Is(err, apperrors.ErrAccountLocked) {
			c.JSON(http.StatusForbidden, gin.H{"error": "account is locked"})
			return
		}
		if errors.Is(err, apperrors.ErrAccountInactive) {
			c.JSON(http.StatusForbidden, gin.H{"error": "account is inactive"})
			return
		}
		log.Printf("refresh token error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
