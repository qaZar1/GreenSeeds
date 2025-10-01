package transport

import (
	"sync"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/service"
)

type Transport struct {
	service *service.Service
	infra   *infrastructure.Infrastructure
	mu      sync.RWMutex
}

func NewTransport(
	repo *repository.Repository,
	cfg models.Config,
	infra *infrastructure.Infrastructure,
) *Transport {
	service := service.NewService(repo, cfg, infra)

	return &Transport{
		service,
		infra,
		sync.RWMutex{},
	}
}
