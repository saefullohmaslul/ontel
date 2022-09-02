package usecases

import (
	"context"
	"errors"
	"net/http"

	"github.com/saefullohmaslul/distributed-tracing/internal/adapters/postgres"
	"github.com/saefullohmaslul/distributed-tracing/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/pkg"
)

type CustomerUsecase interface {
	GetDetailProfile(ctx context.Context, params *models.DetailProfileRequest) (data models.DetailProfileResponse, err error)
}

type CustomerUsecaseImpl struct {
	customerPostgres postgres.CustomerPostgres
}

func NewCustomerUsecase(customerPostgres postgres.CustomerPostgres) CustomerUsecase {
	return &CustomerUsecaseImpl{
		customerPostgres: customerPostgres,
	}
}

func (u *CustomerUsecaseImpl) GetDetailProfile(ctx context.Context, params *models.DetailProfileRequest) (data models.DetailProfileResponse, err error) {
	customer, err := u.customerPostgres.GetCustomer(ctx, &models.CustomerRequest{
		CustomerID: params.CustomerID,
	})

	if err != nil {
		err = pkg.NewHTTPError(http.StatusInternalServerError, err)
		return
	}

	if (customer == models.Customer{}) {
		err = pkg.NewHTTPError(http.StatusNotFound, errors.New("customer not found"))
		return
	}

	// TODO: call endpoint order

	return
}
