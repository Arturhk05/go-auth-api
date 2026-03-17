package handlers

import (
	"errors"
	"log"
	"net/http"

	apperrors "github.com/arturhk05/go-auth-api/internal/errors"
	"github.com/arturhk05/go-auth-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile godoc
// @Summary      Get User Profile
// @Description  Retrieve the profile of the authenticated user
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  models.User
// @Failure      401   {object}  models.ErrorResponse  "Missing or invalid authorization header, or token expired/invalid"
// @Failure      404   {object}  models.ErrorResponse  "User not found"
// @Failure      500   {object}  models.ErrorResponse  "Internal server error"
// @Router       /user/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	user, err := h.userService.GetUserById(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		log.Printf("error getting user profile: user_id=%s, error=%v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
