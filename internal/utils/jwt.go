package utils

import (
	"bca-go-final/internal/types"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("invalid token")

const minSecretKeySize = 8

type Payload struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CompanyId  uuid.UUID `json:"company_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	IsLoggedIn bool      `json:"is_logged_in"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type JWTMaker struct {
	secretKey string
}

type Maker interface {
	CreateToken(userInfo types.User, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func NewPayload(u types.User, duration time.Duration) *Payload {
	payload := &Payload{
		ID:         u.Id,
		Email:      u.Email,
		CompanyId:  u.CompanyId,
		Name:       u.Name,
		Role:       u.RoleId,
		IsLoggedIn: true,
		IssuedAt:   time.Now(),
		ExpiredAt:  time.Now().Add(duration),
	}
	return payload
}

func (maker *JWTMaker) CreateToken(userInfo types.User, duration time.Duration) (string, error) {
	payload := NewPayload(userInfo, duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (payload *Payload) Valid() error {
	if payload.IsLoggedIn && time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func GenerateToken(u types.User) (string, error) {
	secretKey := os.Getenv("SECRET")
	jwtMaker, err := NewJWTMaker(secretKey)
	if err != nil {
		return "", err
	}
	token, err := jwtMaker.CreateToken(u, 60*time.Minute)
	if err != nil {
		return "", err
	}
	return token, nil
}
