package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	coreLog "github.com/Impisigmatus/service_core/log"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
	"github.com/rs/zerolog"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

func BearerAuthMiddleware(infra *infrastructure.Infrastructure, repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(token, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := infra.GetTokenClaims(token[7:]) // убираем Bearer
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, log, ok, err := validateJWT(*claims, repo, r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), coreLog.CtxKey, log)
			ctx = context.WithValue(ctx, UserCtxKey, user)
			log = log.With().Str("username", claims.Username).Logger()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleRequired(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserCtxKey).(models.User)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if user.IsAdmin != nil && !*user.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func validateJWT(claims models.Claims, repo *repository.Repository, r *http.Request) (models.User, zerolog.Logger, bool, error) {
	now := time.Now().Unix()

	// Проверяем, что токен не просрочен
	if claims.ExpiresAt.Unix() < now {
		return models.User{}, zerolog.Logger{}, false, ErrTokenExpired
	}

	// Проверяем, что токен уже активен
	if claims.IssuedAt.Unix() > now {
		return models.User{}, zerolog.Logger{}, false, ErrTokenNotValidYet
	}

	// Проверяем, что токен предназначен нам
	if claims.Audience[0] != "service.green_seeds.api" {
		return models.User{}, zerolog.Logger{}, false, ErrInvalidAudience
	}

	user, err := repo.UsrRepo.CheckUserById(*claims.UserId)
	if err != nil {
		return models.User{}, zerolog.Logger{}, false, ErrInvalidSubject
	}

	if user == (models.User{}) {
		return models.User{}, zerolog.Logger{}, false, ErrInvalidSubject
	}

	log, ok := r.Context().Value(coreLog.CtxKey).(zerolog.Logger)
	if !ok {
		return models.User{}, zerolog.Logger{}, false, ErrInvalidSubject
	}

	log = log.With().Str("username", claims.Username).Logger()

	return user, log, true, nil
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()

		// log.Info().
		// 	Str("method", r.Method).
		// 	Str("url", r.URL.String()).
		// 	Str("remote", r.RemoteAddr).
		// 	Msg("Incoming request")

		next.ServeHTTP(w, r)

		// log.Info().
		// 	Str("url", r.URL.String()).
		// 	Dur("duration", time.Since(start)).
		// 	Msg("Request completed")
	})
}

func WsAuthMiddleware(infra *infrastructure.Infrastructure, repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
}
