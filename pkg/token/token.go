package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"

	"time"
)

// Common errors
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific user
	CreateToken(userID string, email string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid
	VerifyToken(token string) (*Payload, error)
}

// JWTMaker implements the Maker interface using JWT
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new JWT token
func (maker *JWTMaker) CreateToken(userID string, email string, duration time.Duration) (string, error) {
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	claims := jwt.MapClaims{
		"user_id":    userID,
		"email":      email,
		"issued_at":  issuedAt.Unix(),
		"expired_at": expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.secretKey))
}

// VerifyToken verifies if the token is valid
func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	expiredAt, ok := claims["expired_at"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	issuedAt, ok := claims["issued_at"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		UserID:    userID,
		Email:     email,
		IssuedAt:  time.Unix(int64(issuedAt), 0),
		ExpiredAt: time.Unix(int64(expiredAt), 0),
	}

	// Check if token is expired
	if time.Now().After(payload.ExpiredAt) {
		return nil, ErrExpiredToken
	}

	return payload, nil
}
