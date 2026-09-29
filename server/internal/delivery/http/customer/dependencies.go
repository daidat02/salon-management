package customer

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/customer"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.CustomerUsecase
	Handler *CustomerHandler
	Repo    *repo.PostgresCustomerRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	customerRepo := repo.NewPostgresCustomerRepository(dbPool)
	appointmentRepo := repo.NewPostgresAppointmentRepository(dbPool)
	useCase := uc.NewCustomerUsecase(customerRepo, appointmentRepo)
	handler := NewCustomerHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo:    customerRepo,
	}
}
