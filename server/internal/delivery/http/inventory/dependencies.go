package inventory

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/inventory"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.InventoryUsecase
	Handler *InventoryHandler
	Repo    *repo.PostgresInventoryRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	inventoryRepo := repo.NewPostgresInventoryRepository(dbPool)
	serviceRepo := repo.NewPostgresServiceRepository(dbPool)
	useCase := uc.NewInventoryUsecase(inventoryRepo, serviceRepo)
	handler := NewInventoryHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo:    inventoryRepo,
	}
}
