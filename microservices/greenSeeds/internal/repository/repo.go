package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	UsrRepo IUsersRepository
	RepRepo IReportsRepository
	ShfRepo IShiftsRepository
	PlcRepo IPlacementRepository
	RptRepo IReceiptsRepository
	AsnRepo IAssignmentsRepository
	SedRepo ISeedsRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UsrRepo: NewUsersRepository(db),
		RepRepo: NewReportsRepository(db),
		ShfRepo: NewShiftsRepository(db),
		PlcRepo: NewPlacementRepository(db),
		RptRepo: NewReceiptsRepository(db),
		AsnRepo: NewAssignmentsRepository(db),
		SedRepo: NewSeedsRepository(db),
	}
}
