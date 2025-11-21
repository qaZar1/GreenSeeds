package repository

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
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
}

func NewRepository(db *sqlx.DB) *Repository {
	usersRepo := NewUsersRepository(db)

	users, err := usersRepo.CheckAllUsers()
	if err != nil {
		log.Panic().Err(err).Msg("Can not check users in DB")
	}
	if users == nil {
		uspass := "admin"
		admin := true
		usersRepo.AddUser(models.User{
			Username: uspass,
			Password: &uspass,
			IsAdmin:  &admin,
		})
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
	}
}
