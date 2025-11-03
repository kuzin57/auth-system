package handlers

import (
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuzin57/auth-system/internal/services/token"
)

//go:embed templates/registration.html
var registrationTemplate string

type templateData struct {
	Email string
	Token string
}

type RegistrationPageHandler struct {
	tokenService *token.Service
}

func NewRegistrationPageHandler(tokenService *token.Service) *RegistrationPageHandler {
	return &RegistrationPageHandler{
		tokenService: tokenService,
	}
}

func (h *RegistrationPageHandler) Handle(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.String(http.StatusBadRequest, "Token is required")
		return
	}

	claims, err := h.tokenService.ValidateToken(tokenStr)
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	tmpl, err := template.New("registration").Parse(registrationTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse template")
		return
	}

	data := templateData{
		Email: claims.Email,
		Token: tokenStr,
	}

	var html strings.Builder
	if err := tmpl.Execute(&html, data); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render template")
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html.String()))
}
