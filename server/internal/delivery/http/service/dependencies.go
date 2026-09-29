package service

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.ServiceUsecase
	Handler *ServiceHandler
	Repo    *repo.PostgresServiceRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	repo := repo.NewPostgresServiceRepository(dbPool)
	useCase := uc.NewServiceUsecase(repo)
	handler := NewServiceHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo:    repo,
	}
}
