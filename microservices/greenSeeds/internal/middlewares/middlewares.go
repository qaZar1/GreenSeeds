package middlewares

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/infrastructure"
// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
// 	"github.com/rs/zerolog/log"
// )

// func BearerAuthMiddleware(infra *infrastructure.Infrastructure, repo *repository.Repository) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			token := r.Header.Get("Authorization")
// 			if token == "" {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			if !strings.HasPrefix(token, "Bearer ") {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			claims, err := infra.GetTokenClaims(token[7:]) // убираем Bearer
// 			if err != nil {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			if claims == nil {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			ok, err := validateJWT(*claims, repo)
// 			if err != nil {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			if !ok {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func validateJWT(claims models.Claims, repo *repository.Repository) (bool, error) {
// 	now := time.Now().Unix()

// 	// Проверяем, что токен не просрочен (exp)
// 	if claims.ExpiresAt.Unix() < now {
// 		return false, ErrTokenExpired
// 	}

// 	// Проверяем, что токен уже активен
// 	if claims.IssuedAt.Unix() > now {
// 		return false, ErrTokenNotValidYet
// 	}

// 	// Проверяем, что токен предназначен нашему сервису (aud)
// 	if claims.Audience[0] != "service.green_seeds.api" {
// 		return false, ErrInvalidAudience
// 	}

// 	user, err := repo.UsrRepo.CheckUserByUuid(claims.Subject)
// 	if err != nil {
// 		return false, ErrInvalidSubject
// 	}

// 	if user == (models.User{}) {
// 		return false, ErrInvalidSubject
// 	}

// 	if claims.Resources["service.green_seeds.api"].Roles[0] != "service.green_seeds.api:"+*user.Role {
// 		return false, ErrInvalidSubject
// 	}

// 	return true, nil
// }

// func LoggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		log.Info().
// 			Str("method", r.Method).
// 			Str("url", r.URL.String()).
// 			Str("remote", r.RemoteAddr).
// 			Msg("Incoming request")

// 		next.ServeHTTP(w, r)

// 		log.Info().
// 			Str("url", r.URL.String()).
// 			Dur("duration", time.Since(start)).
// 			Msg("Request completed")
// 	})
// }
