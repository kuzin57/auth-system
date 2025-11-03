package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuzin57/auth-system/internal/models"
	"github.com/kuzin57/auth-system/internal/services/users"
)

type SendRegistrationLinkHandler struct {
	service *users.Service
}

func NewSendRegistrationLinkHandler(service *users.Service) *SendRegistrationLinkHandler {
	return &SendRegistrationLinkHandler{service: service}
}

func (h *SendRegistrationLinkHandler) Handle(c *gin.Context) {
	var request models.SendRegistrationLinkRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SendRegistrationLink(c.Request.Context(), request)
	if err != nil {
		log.Println("failed to send registration link", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
