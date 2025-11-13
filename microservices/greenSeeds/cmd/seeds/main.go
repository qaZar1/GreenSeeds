package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/config"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/postgres"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/router"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/ws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"

	_ "gopkg.in/yaml.v3"
)

// @title ITops API
// @version 1.0
// @description API для работы с ITops
// @BasePath /
func main() {
	fmt.Println(time.Now())
	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = "config.yaml"
	}

	cfg := config.MakeConfig(configPath)
	if cfg == (models.Config{}) {
		log.Error().Msg("Invalid config")
		return
	}

	if err := os.MkdirAll(cfg.Server.PathToLog, 0755); err != nil {
		log.Error().Err(err).Msg("Cannot create log directory")
		return
	}

	file, err := os.OpenFile(cfg.Server.PathToLog+"/server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Error().Err(err).Msg("Cannot open log file")
		return
	}
	defer file.Close()

	log.Logger = zerolog.New(file).With().
		Timestamp().
		Logger()

	log.Info().Msg("Logger initialized")

	postgres := sqlx.NewDb(postgres.NewPostgres(
		postgres.Config{
			Hostname: cfg.Database.Host,
			Port:     cfg.Database.Port,
			Database: cfg.Database.Name,
			User:     cfg.Database.User,
			Password: cfg.Database.Pass,
		},
	), "pgx")

	if postgres.DB == nil {
		log.Error().Msg("Invalid postgres connect")
		return
	}

	repo := repository.NewRepository(
		postgres,
	)

	camera := camera.NewCamera(
		cfg.Camera.Name,
		cfg.Camera.InputDevice,
		cfg.Camera.Framerate,
		cfg.Camera.VideoSize,
	)

	ws, err := ws.NewServer(cfg.Serial.Port, cfg.Serial.Baud, camera)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create ws server")
	}
	defer ws.Close()

	router := router.NewRouter(repo, cfg, ws)
	if router == nil {
		log.Error().Msg("Invalid router")
		return
	}

	server := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return
		}
		log.Info().Msg("Service stopped")
	}()
	log.Info().Msg("🚀 Запуск программы: Hortus")
	log.Info().Msg("Service started on " + cfg.Server.Address)

	channel := make(chan os.Signal, 1)
	signal.Notify(channel,
		syscall.SIGABRT,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-channel

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Service shutdown\n\n")
	}
}
