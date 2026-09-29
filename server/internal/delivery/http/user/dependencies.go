package user

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.UserUsecase
	Handler *UserHandler
	Repo *repo.PostgresUserRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	repo := repo.NewPostgresUserRepository(dbPool) // Assuming you have a database connection 'db' available
	useCase := uc.NewCreateUserUsecase(repo)
	handler := NewCreateUserHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo: repo,
	}
}