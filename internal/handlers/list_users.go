package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuzin57/auth-system/internal/services/users"
)

type ListUsersHandler struct {
	service *users.Service
}

func NewListUsersHandler(service *users.Service) *ListUsersHandler {
	return &ListUsersHandler{service: service}
}

func (h *ListUsersHandler) Handle(c *gin.Context) {
	log.Println("request", c.Request)

	response, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		log.Println("failed to list users", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
