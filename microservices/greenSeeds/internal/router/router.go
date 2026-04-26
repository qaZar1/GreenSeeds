package router

import (
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	core "github.com/Impisigmatus/service_core/middlewares"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/camera"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/middlewares"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/ws"
	"github.com/rs/zerolog"
)

func NewRouter(
	repo *repository.Repository,
	cfg models.Config,
	ws *ws.Server,
	logger zerolog.Logger,
	camera *camera.Camera,
	infra *infrastructure.Infrastructure,
) *chi.Mux {
	transport := transport.NewTransport(repo, cfg, infra, ws, camera)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/auth/login", transport.PostApiLoginUser)

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	router.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	router.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	router.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	router.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	// swagger
	// docs.SwaggerInfo.Title = "GreenSeeds API"
	// docs.SwaggerInfo.Description = "API для работы с GreenSeeds"
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8001"
	// docs.SwaggerInfo.BasePath = "/"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.Handle("/swagger/*", httpSwagger.WrapHandler)

	router.HandleFunc("/ws", ws.HandleWS)

	router.Route("/api", func(r chi.Router) {
		r.Use(
			core.RequestID(logger),
			core.ContextLogger(),
			middlewares.BearerAuthMiddleware(infra, repo),
		)

		r.Get("/shifts/getWithoutUser", transport.GetApiShiftsGetWithoutUser)
		r.Put("/shifts/update", transport.PutApiShiftsUpdate)

		r.Get("/assignments/active-tasks/{user_id}", transport.GetApiActiveTasks)
		r.Get("/assignments/task/{id}", transport.GetApiTask)

		r.Get("/users/get/{user_id}", transport.GetApiUserGetUsername)
		r.Put("/users/update", transport.PutApiUpdateUser)
		r.Put("/users/change-password", transport.PutApiChangePassword)

		r.Route("/admin", func(r chi.Router) {
			r.Use(middlewares.RoleRequired("admin"))
			r.Post("/register", transport.PostApiRegisterUser)

			r.Route("/seeds", func(r chi.Router) {
				r.Post("/add", transport.PostApiSeedAdd)
				r.Get("/get", transport.GetApiSeedGet)
				r.Get("/get/{seed}", transport.GetApiSeedGetSeed)
				r.Put("/update", transport.PutApiSeedUpdate)
				r.Delete("/delete/{seed}", transport.DeleteApiSeedDelete)
				r.Get("/getWithBunkers/{seed}", transport.GetApiSeedWithBunkers)
			})

			r.Route("/bunkers", func(r chi.Router) {
				r.Post("/add", transport.PostApiBunkerAdd)
				r.Get("/get", transport.GetApiBunkerGet)
				r.Get("/get/{bunker}", transport.GetApiBunkerGetId)
				r.Put("/update", transport.PutApiBunkerUpdate)
				r.Delete("/delete/{bunker}", transport.DeleteApiBunkerDelete)
				r.Get("/getForPlacement", transport.GetApiBunkerGetForPlacement)
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/get", transport.GetApiCheckAllUsers)
				r.Delete("/delete/{username}", transport.DeleteApiRemoveUser)
			})

			r.Route("/placement", func(r chi.Router) {
				r.Post("/add", transport.PostApiPlacementAdd)
				r.Get("/get", transport.GetApiPlacementGet)
				r.Get("/get/{bunker}", transport.GetApiPlacementGetBunker)
				r.Put("/update", transport.PutApiPlacementUpdate)
				r.Delete("/delete/{bunker}", transport.DeleteApiPlacementDelete)
				r.Put("/fill", transport.PutApiPlacementFill)
			})

			r.Route("/receipts", func(r chi.Router) {
				r.Post("/add", transport.PostApiReceiptsAdd)
				r.Get("/get", transport.GetApiReceiptsGet)
				r.Get("/get/{receipt}", transport.GetApiReceiptsGetReceipt)
				r.Put("/update", transport.PutApiReceiptsUpdate)
				r.Delete("/delete/{receipt}", transport.DeleteApiReceiptsDelete)
			})

			r.Route("/shifts", func(r chi.Router) {
				r.Post("/add", transport.PostApiShiftAdd)
				r.Get("/get", transport.GetApiShiftsGet)
				r.Get("/get/{shift}", transport.GetApiShiftsGetShift)

				r.Delete("/delete/{shift}", transport.DeleteApiShiftsDelete)
			})

			r.Route("/assignments", func(r chi.Router) {
				r.Post("/add", transport.PostApiAssignmentsAdd)
				r.Get("/get", transport.GetApiAssignmentsGet)
				r.Get("/get/{id}", transport.GetApiAssignmentsGetAssignment)
				r.Put("/update", transport.PutApiAssignmentsUpdate)
				r.Delete("/delete/{id}", transport.DeleteApiAssignmentsDelete)
			})

			r.Route("/reports", func(r chi.Router) {
				r.Post("/add", transport.PostApiReportsAdd)
				r.Get("/get", transport.GetApiReports)
				r.Get("/get/{id}", transport.GetApiReportsById)
			})

			r.Route("/logs", func(r chi.Router) {
				r.Get("/get", transport.GetApiLogsGet)
			})

			r.Route("/device-settings", func(r chi.Router) {
				r.Post("/add", transport.PostApiDeviceSettingsAdd)
				r.Get("/get", transport.GetApiDeviceSettingsGet)
				r.Get("/get/{key}", transport.GetApiDeviceSettingsGetKey)
				r.Put("/update", transport.PutApiDeviceSettingsUpdate)
				r.Delete("/delete/{key}", transport.DeleteApiDeviceSettingsDelete)
			})
		})

		r.Route("/calibration", func(r chi.Router) {
			r.Post("/handshake", transport.PostApiCalibrationHandshake)
			r.Post("/photo/{number-of-photo}", transport.PostApiCalibrationPhoto)
			r.Post("/clear", transport.PostApiCalibrationClear)
			r.Post("/calculate", transport.PostApiCalibrationCalc)
			r.Post("/save", transport.PostApiCalibrationSave)
			// r.Get("/stream", transport.GetApiCalibrationStream)
		})
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := "./dist" + r.URL.Path
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		http.ServeFile(w, r, "./dist/index.html")
	})

	return router
}
