package application

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/opencv"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/ws"
)

type App struct {
	repo        *repository.Repository
	infra       *infrastructure.Infrastructure
	validate    *validator.Validate
	cfg         models.Config
	ws          *ws.Server
	camera      *camera.Camera
	calibration map[string]models.Calibration
	calib       *opencv.Calibration
}

func NewApp(
	repo *repository.Repository,
	cfg models.Config,
	infra *infrastructure.Infrastructure,
	ws *ws.Server,
	camera *camera.Camera,
) *App {
	validate := validator.New()

	calibration := opencv.NewCalibration()
	return &App{
		repo:        repo,
		infra:       infra,
		validate:    validate,
		cfg:         cfg,
		ws:          ws,
		camera:      camera,
		calibration: make(map[string]models.Calibration),
		calib:       calibration,
	}
}
