package transport

import (
	"sync"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/application"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/ws"
)

type Transport struct {
	app    *application.App
	infra  *infrastructure.Infrastructure
	mu     sync.RWMutex
	ws     *ws.Server
	camera *camera.Camera
}

func NewTransport(
	repo *repository.Repository,
	cfg models.Config,
	infra *infrastructure.Infrastructure,
	ws *ws.Server,
	camera *camera.Camera,
) *Transport {
	app := application.NewApp(repo, cfg, infra, ws, camera)

	return &Transport{
		app,
		infra,
		sync.RWMutex{},
		ws,
		camera,
	}
}
