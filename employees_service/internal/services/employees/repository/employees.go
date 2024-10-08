package employees_repository_repository

import "github.com/Sleeps17/events-planning-service-backend/employees_service/internal/repository"

type service struct {
	repo repository.EmployeeRepository
}

func New(repo repository.EmployeeRepository) *service {
	return &service{repo: repo}
}
