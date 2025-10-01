package infrastructure

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

const resource = "service.green_seeds.api"

func (infra *Infrastructure) GetClaims(uuid string, role string) models.Claims {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Duration(infra.ExpiresIn) * time.Second)

	return models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Audience: jwt.ClaimStrings{
				resource,
			},
		},

		Subject:   uuid,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		ExpiresIn: infra.ExpiresIn,
		Type:      "Bearer",
		Resources: map[string]models.Roles{
			resource: {
				Roles: []string{
					resource + ":" + role,
				},
			},
		},
	}
}

func (infra *Infrastructure) GetSignedToken(claims models.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(infra.secret)
}

func (infra *Infrastructure) GetTokenClaims(auth string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(
		auth,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return infra.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, TokenError
	}

	return claims, nil
}
