package appointment

import (
	repo "github.com/daidat02/server/internal/infrastructure/postgres"
	uc "github.com/daidat02/server/internal/usecase/appointment"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	UseCase *uc.AppointmentUseCase
	Handler *AppointmentHandler
	Repo    *repo.PostgresAppointmentRepository
}

func NewDependencies(dbPool *pgxpool.Pool) *Dependencies {
	appointmentRepo := repo.NewPostgresAppointmentRepository(dbPool)
	customerRepo := repo.NewPostgresCustomerRepository(dbPool)
	serviceRepo := repo.NewPostgresServiceRepository(dbPool)
	useCase := uc.NewAppointmentUseCase(appointmentRepo, customerRepo, serviceRepo)
	handler := NewAppointmentHandler(useCase)

	return &Dependencies{
		UseCase: useCase,
		Handler: handler,
		Repo:    appointmentRepo,
	}
}
