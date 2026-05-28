package transport

import (
	"sync"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/application"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/ws"
)

type Transport struct {
	Assignments    application.IAssignmentsApp
	Bunkers        application.IBunkersApp
	DeviceSettings application.IDeviceSettingsApp
	Placements     application.IPlacementsApp
	Recipes        application.IRecipesApp
	Reports        application.IReportsApp
	Seeds          application.ISeedsApp
	Shifts         application.IShiftsApp
	Users          application.IUsersApp
	Calibration    application.ICalibrationApp
	Logs           application.ILogsApp

	infra  *infrastructure.Infrastructure
	mu     sync.RWMutex
	ws     *ws.Server
	camera camera.ICamera
}

func NewTransport(
	repo *repository.Repository,
	cfg models.Config,
	infra *infrastructure.Infrastructure,
	ws *ws.Server,
	camera camera.ICamera,
	client *device.DeviceClient,
) *Transport {
	app := application.NewApp(repo, cfg, infra, ws, camera, client)
	return &Transport{
		Assignments:    app,
		Bunkers:        app,
		DeviceSettings: app,
		Placements:     app,
		Recipes:        app,
		Reports:        app,
		Seeds:          app,
		Shifts:         app,
		Users:          app,
		Calibration:    app,
		Logs:           app,
		infra:          infra,
		mu:             sync.RWMutex{},
		ws:             ws,
		camera:         camera,
	}
}
