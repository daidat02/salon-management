package product

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/product"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.ProductUsecase
	Handler *ProductHandler
	Repo    *repo.PostgresProductRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	repo := repo.NewPostgresProductRepository(dbPool)
	useCase := uc.NewProductUsecase(repo)
	handler := NewProductHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo:    repo,
	}
}
