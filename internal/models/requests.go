package models

type SendRegistrationLinkRequest struct {
	Email string `json:"email"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
