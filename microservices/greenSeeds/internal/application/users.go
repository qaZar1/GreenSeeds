package application

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func (app *App) RegisterUser(user models.User) (int, error) {
	if err := app.validate.Struct(user); err != nil {
		return http.StatusBadRequest, ErrValidateStruct
	}

	if user.Password == nil {
		user.Password = &app.cfg.Auth.DefaultPassword
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

	ok, err := app.repo.UsrRepo.AddUser(newUser)
	if err != nil {
		return http.StatusInternalServerError, ErrInvalidAddUser
	}

	if !ok {
		return http.StatusBadRequest, ErrUserAlreadyExists
	}

	return http.StatusNoContent, nil
}

func (app *App) LoginUser(user models.User) (*models.TokenResponse, int, error) {
	if err := app.validate.Struct(user); err != nil {
		return nil, http.StatusBadRequest, ErrValidateStruct
	}

	checkedUser, err := app.repo.UsrRepo.CheckUserByUsernameWithPwd(user.Username)
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

	claims := app.infra.GetClaims(*checkedUser.Id, checkedUser.Username, role, *checkedUser.FullName)

	signed, err := app.infra.GetSignedToken(claims)
	if err != nil {
		return nil, http.StatusInternalServerError, ErrInvalidGenerateToken
	}

	return &models.TokenResponse{
		AccessToken: signed,
		ExpiresIn:   app.cfg.JWT.ExpiresIn,
		TokenType:   claims.Type,
	}, http.StatusOK, nil
}

func (app *App) GetUserById(id string) (*models.User, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := app.repo.UsrRepo.CheckUserById(idInt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (app *App) CheckAllUsers() ([]models.User, error) {
	allUsers, err := app.repo.UsrRepo.CheckAllUsers()
	if err != nil {
		return nil, err
	}

	return allUsers, nil
}

func (app *App) Update(user models.User) (bool, error) {
	if err := app.validate.Struct(user); err != nil {
		return false, ErrValidateStruct
	}

	result, err := app.repo.UsrRepo.Update(user)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (app *App) ChangePassword(user models.UpdatePassword) (bool, error) {
	if err := app.validate.Struct(user); err != nil {
		return false, err
	}

	userData, err := app.repo.UsrRepo.CheckUserById(user.Id)
	if err != nil {
		return false, err
	}

	checkedUser, err := app.repo.UsrRepo.CheckUserByUsernameWithPwd(userData.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
		}
		return false, err
	}

	if user.OldPassword == nil && user.NewPassword == nil {
		user.NewPassword = &app.cfg.Auth.DefaultPassword
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

	result, err := app.repo.UsrRepo.UpdatePassword(user)

	return result, err
}

func (app *App) RemoveUser(username string) (bool, error) {
	return app.repo.UsrRepo.Delete(username)
}
