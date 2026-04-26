package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims

	Username  string    `json:"username"`
	Role      string    `json:"role"`
	UserId    *int64    `json:"user_id"`
	FullName  string    `json:"full_name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	ExpiresIn int       `json:"expires_in"`
	Type      string    `json:"type"`

	Resources map[string]string `json:"resource_access"`
}

type JWT struct {
	Exp         int    `json:"exp"`
	Customer    string `json:"customer"`
	Inn         string `json:"inn"`
	Product     string `json:"product"`
	Channels    string `json:"channels"`
	Delivers    string `json:"delivers"`
	Dispatchers string `json:"dispatchers"`
	ValidFrom   string `json:"validFrom"`
	Valid       bool   `json:"valid"`
}
