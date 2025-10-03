package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/config"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/postgres"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/router"
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

	if postgres == nil {
		log.Error().Msg("Invalid postgres connect")
		return
	}

	repo := repository.NewRepository(
		postgres,
	)

	router := router.NewRouter(repo, cfg)
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
