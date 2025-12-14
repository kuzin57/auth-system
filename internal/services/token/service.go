package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kuzin57/auth-system/internal/config"
)

type Service struct {
	secretKey []byte
}

type Claims struct {
	jwt.RegisteredClaims

	Email string `json:"email"`
}

func NewService(secrets *config.Secrets) *Service {
	return &Service{
		secretKey: []byte(secrets.JWT.SecretKey),
	}
}

func (s *Service) GenerateToken(email string, expirationTime time.Duration) (string, error) {
	if len(s.secretKey) == 0 {
		return "", errors.New("secret key is not set")
	}

	expiration := time.Now().Add(expirationTime)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	if len(s.secretKey) == 0 {
		return nil, errors.New("secret key is not set")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
