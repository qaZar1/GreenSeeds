package infrastructure

import "github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"

type Infrastructure struct {
	ExpiresIn int

	secret []byte

	resources map[string]models.Roles
}

func New(expiresIn int, cfg models.Config) *Infrastructure {
	return &Infrastructure{
		secret:    []byte(cfg.JWT.Secret),
		ExpiresIn: expiresIn,
		resources: map[string]models.Roles{},
	}
}
