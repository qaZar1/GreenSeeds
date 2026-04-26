package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/sqlite"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	UsrRepo  IUsersRepository
	RepRepo  IReportsRepository
	ShfRepo  IShiftsRepository
	PlcRepo  IPlacementRepository
	RptRepo  IReceiptsRepository
	AsnRepo  IAssignmentsRepository
	SeedRepo ISeedsRepository
	BunkRepo IBunkersRepository
	LogsRepo ILogsRepository
	CalRepo  ICalibrationsRepository
	DevSet   IDeviceSettingsRepository

	SQLite *sqlite.SQLite
}

func NewRepository(db *sqlx.DB, sqlite *sqlite.SQLite) *Repository {
	usersRepo := NewUsersRepository(db)

	users, err := usersRepo.CheckAllUsers()
	if err != nil {
		log.Panic().Err(err).Msg("Can not check users in DB")
	}
	if users == nil {
		pass := "admin"

		newPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
		if err != nil {
			log.Panic().Err(err).Msg("Can not generate password hash")
		}

		tmp := string(newPass)
		admin := true

		newUser := models.User{
			Username: pass,
			Password: &tmp,
			FullName: &pass,
			IsAdmin:  &admin,
		}
		// todo сделать проверку на ошибку
		// todo пароль должен быть хэширован

		if _, err := usersRepo.AddUser(newUser); err != nil {
			log.Panic().Err(err).Msg("Can not add user to DB")
		}
	}

	return &Repository{
		UsrRepo:  usersRepo,
		RepRepo:  NewReportsRepository(db),
		ShfRepo:  NewShiftsRepository(db),
		PlcRepo:  NewPlacementRepository(db),
		RptRepo:  NewReceiptsRepository(db),
		AsnRepo:  NewAssignmentsRepository(db),
		SeedRepo: NewSeedsRepository(db),
		BunkRepo: NewBunkersRepository(db),
		LogsRepo: NewLogsRepository(db),
		CalRepo:  NewCalibrationsRepository(db),
		DevSet:   NewDeviceSettingsRepository(db),

		SQLite: sqlite,
	}
}
