package category

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/category"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.CategoryUsecase
	Handler *CategoryHandler
	Repo    *repo.PostgresCategoryRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	repo := repo.NewPostgresCategoryRepository(dbPool)
	useCase := uc.NewCategoryUsecase(repo)
	handler := NewCategoryHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo:    repo,
	}
}
