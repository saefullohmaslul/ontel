package usecases

import (
	"context"
	"errors"
	"net/http"

	httpclient "github.com/saefullohmaslul/distributed-tracing/customer-service/internal/adapters/http"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/adapters/postgres"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/pkg"
)

type CustomerUsecase interface {
	GetDetailProfile(ctx context.Context, params *models.DetailProfileRequest) (data models.DetailProfileResponse, err error)
}

type CustomerUsecaseImpl struct {
	customerPostgres postgres.CustomerPostgres
	orderHttp        httpclient.OrderHttpClient
}

func NewCustomerUsecase(customerPostgres postgres.CustomerPostgres, orderHttp httpclient.OrderHttpClient) CustomerUsecase {
	return &CustomerUsecaseImpl{
		customerPostgres: customerPostgres,
		orderHttp:        orderHttp,
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

	orders, err := u.orderHttp.GetDetailOrders(ctx, &models.OrdersDetailRequest{
		CustomerID: params.CustomerID,
	})

	if err != nil {
		err = pkg.NewHTTPError(http.StatusNotFound, err)
		return
	}

	data.Customer = customer
	data.Orders = orders

	return
}
