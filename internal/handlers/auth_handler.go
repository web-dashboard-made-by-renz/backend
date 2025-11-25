package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/web-dashboard-made-by-renz/backend/internal/models"
	"github.com/web-dashboard-made-by-renz/backend/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	// Get user info from context (set by middleware)
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"user": models.User{
			Username: username.(string),
			Role:     role.(string),
		},
	})
}
