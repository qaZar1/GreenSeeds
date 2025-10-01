package service

// import (
// 	"database/sql"
// 	"net/http"

// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
// 	bcrypt "golang.org/x/crypto/bcrypt"
// )

// func (s *Service) RegisterUser(user models.User) (int, error) {
// 	if err := s.validate.Struct(user); err != nil {
// 		return http.StatusBadRequest, ErrValidateStruct
// 	}

// 	newPass, err := bcrypt.GenerateFromPassword([]byte(*user.PasswordHash), 10)
// 	if err != nil {
// 		return http.StatusInternalServerError, ErrInvalidGeneratePassHash
// 	}

// 	tmp := string(newPass)

// 	newUser := models.User{
// 		FullName:     user.FullName,
// 		PasswordHash: &tmp,
// 		Role:         user.Role,
// 	}

// 	ok, err := s.repo.UsrRepo.AddUser(newUser)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return http.StatusConflict, ErrUserAlreadyExists
// 		}
// 		return http.StatusInternalServerError, ErrInvalidAddUser
// 	}

// 	if !ok {
// 		return http.StatusBadRequest, ErrUserAlreadyExists
// 	}

// 	return http.StatusNoContent, nil
// }

// func (s *Service) LoginUser(user models.User) (*models.TokenResponse, int, error) {
// 	if err := s.validate.Struct(user); err != nil {
// 		return nil, http.StatusBadRequest, ErrValidateStruct
// 	}

// 	checkedUser, err := s.repo.UsrRepo.CheckUserByFullName(user.FullName)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, http.StatusNotFound, ErrUserNotFound
// 		}
// 		return nil, http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	if bcrypt.CompareHashAndPassword([]byte(*checkedUser.PasswordHash), []byte(*user.PasswordHash)) != nil {
// 		return nil, http.StatusUnauthorized, ErrInvalidUsernameOrPassword
// 	}

// 	role, err := s.repo.UsrRepo.CheckRolesById(*checkedUser.UUID)
// 	if err != nil {
// 		return nil, http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	claims := s.infra.GetClaims(*checkedUser.UUID, role)

// 	signed, err := s.infra.GetSignedToken(claims)
// 	if err != nil {
// 		return nil, http.StatusInternalServerError, ErrInvalidGenerateToken
// 	}

// 	return &models.TokenResponse{
// 		AccessToken: signed,
// 		ExpiresIn:   s.cfg.JWT.ExpiresIn,
// 	}, http.StatusOK, nil
// }

// func (s *Service) CheckUserByUuid(uuid string) (*models.User, int, error) {
// 	user, err := s.repo.UsrRepo.CheckUserByUuid(uuid)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, http.StatusNotFound, ErrUserNotFound
// 		}

// 		return nil, http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	return &user, http.StatusOK, nil
// }

// func (s *Service) CheckRolesById(uuid string) (string, int, error) {
// 	role, err := s.repo.UsrRepo.CheckRolesById(uuid)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", http.StatusNotFound, ErrUserNotFound
// 		}

// 		return "", http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	return role, http.StatusOK, nil
// }

// func (s *Service) CheckAllUsers() ([]models.User, int, error) {
// 	allUsers, err := s.repo.UsrRepo.CheckAllUsers()
// 	if err != nil {
// 		return nil, http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	return allUsers, http.StatusOK, nil
// }

// func (s *Service) UpdateRole(updateRole models.UpdateRole) (models.UpdateRole, int, error) {
// 	if err := s.validate.Struct(updateRole); err != nil {
// 		return models.UpdateRole{}, http.StatusBadRequest, ErrValidateStruct
// 	}

// 	result, err := s.repo.UsrRepo.UpdateRole(updateRole)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return models.UpdateRole{}, http.StatusNotFound, ErrUserNotFound
// 		}
// 		return models.UpdateRole{}, http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	return result, http.StatusOK, nil
// }

// func (s *Service) ChangePassword(updatePassword models.UpdatePassword) (int, error) {
// 	if err := s.validate.Struct(updatePassword); err != nil {
// 		return http.StatusBadRequest, ErrValidateStruct
// 	}

// 	user, err := s.repo.UsrRepo.CheckUserByUuid(updatePassword.UUID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return http.StatusNotFound, ErrUserNotFound
// 		}
// 		return http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	if bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(*updatePassword.OldPassword)) != nil {
// 		return http.StatusUnauthorized, ErrInvalidUsernameOrPassword
// 	}

// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(updatePassword.NewPassword), 10)
// 	if err != nil {
// 		return http.StatusInternalServerError, ErrInvalidGeneratePassHash
// 	}

// 	updatePassword.NewPassword = string(passwordHash)

// 	result, err := s.repo.UsrRepo.UpdatePassword(updatePassword)
// 	if err != nil {
// 		return http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	if !result {
// 		return http.StatusInternalServerError, ErrInvalidUpdateData
// 	}

// 	return http.StatusOK, nil
// }

// func (s *Service) RemoveUser(removeUser models.RemoveUser) (int, error) {
// 	if err := s.validate.Struct(removeUser); err != nil {
// 		return http.StatusBadRequest, ErrValidateStruct
// 	}

// 	for _, uuid := range removeUser.UUID {
// 		result, err := s.repo.UsrRepo.RemoveUser(uuid)
// 		if err != nil {
// 			return http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 		}

// 		if !result {
// 			return http.StatusInternalServerError, ErrInvalidUpdateData
// 		}
// 	}

// 	return http.StatusOK, nil
// }

// func (s *Service) ResetPassword(uuid string) (int, error) {
// 	user, err := s.repo.UsrRepo.CheckUserByUuid(uuid)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return http.StatusNotFound, ErrUserNotFound
// 		}
// 		return http.StatusInternalServerError, ErrInvalidGetDataFromDB
// 	}

// 	if user == (models.User{}) {
// 		return http.StatusNotFound, ErrUserNotFound
// 	}

// 	newPass, err := bcrypt.GenerateFromPassword([]byte(s.cfg.Auth.DefaultPassword), 10)
// 	if err != nil {
// 		return http.StatusInternalServerError, ErrInvalidGeneratePassHash
// 	}

// 	update := models.UpdatePassword{
// 		UUID:        *user.UUID,
// 		NewPassword: string(newPass),
// 	}

// 	ok, err := s.repo.UsrRepo.UpdatePassword(update)
// 	if err != nil {
// 		return http.StatusInternalServerError, ErrInvalidUpdateData
// 	}

// 	if !ok {
// 		return http.StatusInternalServerError, ErrInvalidUpdateData
// 	}

// 	return http.StatusOK, nil
// }
