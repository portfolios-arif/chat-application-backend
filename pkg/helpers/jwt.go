package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type JWTHelper struct {
	secretKey []byte
	exp       time.Duration
}

func NewJWTHelper(secretKey []byte, exp time.Duration) *JWTHelper {
	return &JWTHelper{
		secretKey: secretKey,
		exp:       exp,
	}
}

func (h *JWTHelper) GenerateToken(userId string) (string, error) {
	claims := &JWTClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *JWTHelper) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return h.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token claims")
	}

	return claims, nil
}

func (h *JWTHelper) ExtractBearerToken(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("Invalid authorization format")
	}

	return authHeader[7:], nil
}
