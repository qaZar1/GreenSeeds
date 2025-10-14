package router

import (
	"net/http"
	"net/http/pprof"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/docs"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/middlewares"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(repo *repository.Repository, cfg models.Config) *chi.Mux {
	infra := infrastructure.New(cfg.JWT.ExpiresIn, cfg)
	transport := transport.NewTransport(repo, cfg, infra)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/api/login", transport.PostApiLoginUser)

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// swagger
	docs.SwaggerInfo.Title = "GreenSeeds API"
	docs.SwaggerInfo.Description = "API для работы с GreenSeeds"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8001"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.Handle("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/api", func(r chi.Router) {
		r.Use(middlewares.BearerAuthMiddleware(infra, repo))
		r.Post("/register", transport.PostApiRegisterUser)

		r.Route("/seeds", func(r chi.Router) {
			r.Post("/add", transport.PostApiSeedAdd)
			r.Get("/get", transport.GetApiSeedGet)
			r.Get("/get/{seed}", transport.GetApiSeedGetId)
			r.Put("/update", transport.PutApiSeedUpdate)
			r.Delete("/delete/{seed}", transport.DeleteApiSeedDelete)
		})

		r.Route("/bunkers", func(r chi.Router) {
			r.Post("/add", transport.PostApiBunkerAdd)
			r.Get("/get", transport.GetApiBunkerGet)
			r.Get("/get/{bunker}", transport.GetApiBunkerGetId)
			r.Put("/update", transport.PutApiBunkerUpdate)
			r.Delete("/delete/{bunker}", transport.DeleteApiBunkerDelete)
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/get", transport.GetApiCheckAllUsers)
			r.Get("/get/{username}", transport.GetApiUserGetUsername)
			r.Put("/update", transport.PutApiUpdateUser)
			r.Put("/change-password", transport.PutApiChangePassword)
			r.Delete("/delete/{username}", transport.DeleteApiRemoveUser)
		})

		r.Route("/placement", func(r chi.Router) {
			r.Post("/add", transport.PostApiPlacementAdd)
			r.Get("/get", transport.GetApiPlacementGet)
			r.Get("/get/{bunker}", transport.GetApiPlacementGetBunker)
			r.Put("/update", transport.PutApiPlacementUpdate)
			r.Delete("/delete/{bunker}", transport.DeleteApiPlacementDelete)
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
			r.Put("/update", transport.PutApiShiftsUpdate)
			r.Delete("/delete/{shift}", transport.DeleteApiShiftsDelete)
		})

		// 	r.Get("/checkByUuid/{uuid}", transport.GetApiCheckUserByUuidUuid)
		// 	r.Get("/checkRoles/{uuid}", transport.GetApiCheckRolesUuid)
		// 	r.Get("/checkAll", transport.GetApiCheckAllUsers)
		// 	r.Put("/updateRole", transport.PutApiUpdateRole)
		// 	r.Put("/change-password", transport.PutApiChangePassword)
		// 	r.Put("/reset-password/{uuid}", transport.PutApiResetPassword)
		// 	r.Delete("/removeUser", transport.DeleteApiRemoveUser)
		// })
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := "./build" + r.URL.Path
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		http.ServeFile(w, r, "./build/index.html")
	})

	return router
}
