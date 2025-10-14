package service

import (
	"database/sql"
	"net/http"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(user models.User) (int, error) {
	if err := s.validate.Struct(user); err != nil {
		return http.StatusBadRequest, ErrValidateStruct
	}

	if user.Password == nil {
		user.Password = &s.cfg.Auth.DefaultPassword
	}

	newPass, err := bcrypt.GenerateFromPassword([]byte(*user.Password), 10)
	if err != nil {
		return http.StatusInternalServerError, ErrInvalidGeneratePassHash
	}

	tmp := string(newPass)

	newUser := models.User{
		Username: user.Username,
		Password: &tmp,
		FullName: user.FullName,
		IsAdmin:  user.IsAdmin,
	}

	ok, err := s.repo.UsrRepo.AddUser(newUser)
	if err != nil {
		return http.StatusInternalServerError, ErrInvalidAddUser
	}

	if !ok {
		return http.StatusBadRequest, ErrUserAlreadyExists
	}

	return http.StatusNoContent, nil
}

func (s *Service) LoginUser(user models.User) (*models.TokenResponse, int, error) {
	if err := s.validate.Struct(user); err != nil {
		return nil, http.StatusBadRequest, ErrValidateStruct
	}

	checkedUser, err := s.repo.UsrRepo.CheckUserByUsernameWithPwd(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, ErrUserNotFound
		}
		return nil, http.StatusInternalServerError, ErrInvalidGetDataFromDB
	}

	if bcrypt.CompareHashAndPassword([]byte(*checkedUser.Password), []byte(*user.Password)) != nil {
		return nil, http.StatusUnauthorized, ErrInvalidUsernameOrPassword
	}

	var role string
	if *checkedUser.IsAdmin {
		role = "admin"
	} else {
		role = "operator"
	}

	claims := s.infra.GetClaims(checkedUser.Username, role)

	signed, err := s.infra.GetSignedToken(claims)
	if err != nil {
		return nil, http.StatusInternalServerError, ErrInvalidGenerateToken
	}

	return &models.TokenResponse{
		AccessToken: signed,
		ExpiresIn:   s.cfg.JWT.ExpiresIn,
		TokenType:   claims.Type,
	}, http.StatusOK, nil
}

func (s *Service) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.repo.UsrRepo.CheckUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) CheckAllUsers() ([]models.User, error) {
	allUsers, err := s.repo.UsrRepo.CheckAllUsers()
	if err != nil {
		return nil, err
	}

	return allUsers, nil
}

func (s *Service) Update(user models.User) (bool, error) {
	if err := s.validate.Struct(user); err != nil {
		return false, ErrValidateStruct
	}

	result, err := s.repo.UsrRepo.Update(user)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *Service) ChangePassword(user models.UpdatePassword) (bool, error) {
	if err := s.validate.Struct(user); err != nil {
		return false, err
	}

	checkedUser, err := s.repo.UsrRepo.CheckUserByUsernameWithPwd(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		return false, err
	}

	if user.OldPassword == nil && user.NewPassword == nil {
		user.NewPassword = &s.cfg.Auth.DefaultPassword
	} else {
		if bcrypt.CompareHashAndPassword([]byte(*checkedUser.Password), []byte(*user.OldPassword)) != nil {
			return false, ErrInvalidUsernameOrPassword
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*user.NewPassword), 10)
	if err != nil {
		return false, err
	}

	hashStr := string(passwordHash)
	user.NewPassword = &hashStr

	result, err := s.repo.UsrRepo.UpdatePassword(user)

	return result, err
}

func (s *Service) RemoveUser(username string) (bool, error) {
	return s.repo.UsrRepo.Delete(username)
}
