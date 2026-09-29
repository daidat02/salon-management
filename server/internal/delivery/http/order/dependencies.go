package order

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.OrderUsecase
	Handler *OrderHandler
	Repo    *repo.PostgresOrderRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	orderRepo := repo.NewPostgresOrderRepository(dbPool)
	serviceRepo := repo.NewPostgresServiceRepository(dbPool)
	useCase := uc.NewOrderUsecase(orderRepo, serviceRepo)
	handler := NewOrderHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo:    orderRepo,
	}
}
