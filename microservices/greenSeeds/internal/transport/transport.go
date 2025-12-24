package transport

import (
	"sync"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/service"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/ws"
)

type Transport struct {
	service *service.Service
	infra   *infrastructure.Infrastructure
	mu      sync.RWMutex
	ws      *ws.Server
	camera  *camera.Camera
}

func NewTransport(
	repo *repository.Repository,
	cfg models.Config,
	infra *infrastructure.Infrastructure,
	ws *ws.Server,
	camera *camera.Camera,
) *Transport {
	service := service.NewService(repo, cfg, infra, ws, camera)

	return &Transport{
		service,
		infra,
		sync.RWMutex{},
		ws,
		camera,
	}
}
