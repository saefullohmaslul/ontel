package postgres

import (
	"context"

	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/pkg"
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
	ctx, span := pkg.NewSpan(ctx, "CustomerPostgresImpl.GetCustomer", nil)
	defer span.End()

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
