package staff

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/staff"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.StaffUsecase
	Handler *StaffHandler
	Repo *repo.PostgresStaffRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	repo := repo.NewPostgresStaffRepository(dbPool)
	useCase := uc.NewStaffUsecase(repo)
	handler := NewStaffHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo: repo,
	}
}