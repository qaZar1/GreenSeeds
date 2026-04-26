package ws

import (
	"errors"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/repository"
)


func WsAuthMiddleware() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(s *Server, client *Client, req models.WSRequest) {
			if !client.IsAuth{
				client.Send <- errResponse(req.Type, errors.New("Unauthorized"))
				return
			}
			
			next(s, client, req)
		}
	}
}

func validateJWT(claims models.Claims, repo *repository.Repository) error {
	now := time.Now().Unix()

	// Проверяем, что токен не просрочен
	if claims.ExpiresAt.Unix() < now {
		return ErrTokenExpired
	}

	// Проверяем, что токен уже активен
	if claims.IssuedAt.Unix() > now {
		return ErrTokenNotValidYet
	}

	// Проверяем, что токен предназначен нам
	if claims.Audience[0] != "service.green_seeds.api" {
		return ErrInvalidAudience
	}

	user, err := repo.UsrRepo.CheckUserById(*claims.UserId)
	if err != nil {
		return ErrInvalidSubject
	}

	if user == (models.User{}) {
		return ErrInvalidSubject
	}

	if *user.IsAdmin {
		return errors.New("Admin has not permission")
	}

	return nil
}