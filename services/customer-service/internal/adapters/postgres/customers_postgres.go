package postgres

import (
	"context"

	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/models"
)

type CustomerPostgres interface {
	GetCustomer(ctx context.Context, params *models.CustomerRequest) (customer models.Customer, err error)
}

type CustomerPostgresImpl struct {
	db *Database
}

func NewCustomerPostgres(db *Database) CustomerPostgres {
	return &CustomerPostgresImpl{
		db: db,
	}
}

func (p *CustomerPostgresImpl) GetCustomer(ctx context.Context, params *models.CustomerRequest) (customer models.Customer, err error) {
	err = p.db.Table("customers").
		Select(
			`
				customer_id,
				name,
				email,
				status
			`,
		).
		Where("customer_id = ?", params.CustomerID).
		WithContext(ctx).
		Find(&customer).
		Error
	return
}
