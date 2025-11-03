package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuzin57/auth-system/internal/models"
	"github.com/kuzin57/auth-system/internal/services/token"
	"github.com/kuzin57/auth-system/internal/services/users"
)

type RegisterHandler struct {
	service      *users.Service
	tokenService *token.Service
}

func NewRegisterHandler(service *users.Service, tokenService *token.Service) *RegisterHandler {
	return &RegisterHandler{
		service:      service,
		tokenService: tokenService,
	}
}

func (h *RegisterHandler) Handle(c *gin.Context) {
	var request models.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := h.tokenService.ValidateToken(request.Token)
	if err != nil {
		log.Println("failed to validate token", err)

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if claims.Email != request.Email {
		log.Println("token email does not match request email")

		c.JSON(http.StatusForbidden, gin.H{"error": "Token email does not match request email"})
		return
	}

	userID, err := h.service.RegisterUser(c.Request.Context(), request)
	if err != nil {
		log.Println("failed to register user", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}
