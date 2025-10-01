package service

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
)

type Service struct {
	repo     *repository.Repository
	infra    *infrastructure.Infrastructure
	validate *validator.Validate
	cfg      models.Config
}

func NewService(repo *repository.Repository, cfg models.Config, infra *infrastructure.Infrastructure) *Service {
	validate := validator.New()
	return &Service{
		repo:     repo,
		infra:    infra,
		validate: validate,
		cfg:      cfg,
	}
}
