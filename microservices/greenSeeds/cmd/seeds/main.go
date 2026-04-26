package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/config"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/logger"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/logger/writer"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/opencv"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/postgres"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/router"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/sqlite"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/ws"
	zerolog "github.com/rs/zerolog"

	_ "gopkg.in/yaml.v3"
)

const logo = `
 _   _   ___   ____   _______   _   _   ______
| | | | / _ \ |  _ \ |_______| | | | | / ____/
| |_| || | | || |_) |   | |    | | | | | (___
|  _  || | | ||  _ <    | |    | | | |  \___ \
| | | || |_| || | \ \   | |    | |_| |  ____) |
|_| |_| \___/ |_|  \_\  |_|     \___/  /_____/
`

// @title GreenSeeds API
// @version 1.0
// @description API для работы c GreenSeeds
// @BasePath /
func main() {
	fmt.Println(logo)
	log := logger.New(zerolog.DebugLevel)

	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		configPath = "config.yaml"
	}

	cfg := config.MakeConfig(configPath)
	if cfg == (models.Config{}) {
		panic("Invalid config")
	}

	postgres := sqlx.NewDb(postgres.NewPostgres(
		postgres.Config{
			Hostname: cfg.Database.Host,
			Port:     cfg.Database.Port,
			Database: cfg.Database.Name,
			User:     cfg.Database.User,
			Password: cfg.Database.Pass,
		},
	), "pgx")

	sqlite := sqlite.NewSQLiteClient(cfg)

	db := writer.NewDbWriter(postgres)
	multi := zerolog.MultiLevelWriter(os.Stdout, db)
	log = log.Output(multi)

	repo := repository.NewRepository(
		postgres,
		sqlite,
	)

	opencv := opencv.NewCounting()
	ctx := context.Background()

	client := device.NewClient(ctx, cfg.Serial.Port, cfg.Serial.Baud, log)

	camera := camera.NewCamera(
		cfg.Camera.Name,
		cfg.Camera.InputDevice,
		cfg.Camera.Framerate,
		cfg.Camera.VideoSize,
	)

	infra := infrastructure.New(cfg.JWT.ExpiresIn, cfg)
	ws, err := ws.NewServer(
		client,
		repo,
		cfg.API.URL,
		log,
		infra,
		camera,
		&opencv,
	)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create ws server")
	}
	defer ws.Close()

	router := router.NewRouter(repo, cfg, ws, log, camera, infra)
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

// // TODO: debug
// func ProcessDataset(root string) error {
// 	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		// пропускаем папки
// 		if d.IsDir() {
// 			return nil
// 		}

// 		ext := strings.ToLower(filepath.Ext(path))
// 		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
// 			return nil
// 		}

// 		// относительный путь
// 		rel, _ := filepath.Rel(root, path)

// 		base := strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel))

// 		// папка вывода
// 		outDir := filepath.Join("out", filepath.Dir(rel), base)
// 		_ = os.MkdirAll(outDir, 0755)

// 		opencv := opencv.NewCounting()
// 		total := opencv.Counter(path, outDir)

// 		fmt.Printf("/%s - total seeds %d\n", rel, total)
// 		return nil
// 	})
// }
