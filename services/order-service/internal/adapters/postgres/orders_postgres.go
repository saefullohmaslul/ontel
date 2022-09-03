package postgres

import (
	"context"

	"github.com/saefullohmaslul/distributed-tracing/order-service/internal/models"
)

type OrderPostgres interface {
	GetOrders(ctx context.Context, params *models.OrdersRequest) (orders []models.Order, err error)
	GetOrderItems(ctx context.Context, params *models.OrderItemsRequest) (orderItems []models.OrderItem, err error)
}

type OrderPostgresImpl struct {
	db *Database
}

func NewOrderPostgres(db *Database) OrderPostgres {
	return &OrderPostgresImpl{
		db: db,
	}
}

func (p *OrderPostgresImpl) GetOrders(ctx context.Context, params *models.OrdersRequest) (orders []models.Order, err error) {
	err = p.db.Table("orders").
		Select(
			`
				order_id,
				customer_id,
				order_no,
				grand_total,
				status,
				created_at,
				updated_at
			`,
		).
		Where("customer_id = ?", params.CustomerID).
		WithContext(ctx).
		Find(&orders).
		Error
	return
}

func (p *OrderPostgresImpl) GetOrderItems(ctx context.Context, params *models.OrderItemsRequest) (orderItems []models.OrderItem, err error) {
	err = p.db.Table("order_items").
		Select(
			`
				order_item_id,	
				order_id,
				sku,
				qty,
				amount,
				created_at,
				updated_at
			`,
		).
		Where("order_id = ?", params.OrderID).
		WithContext(ctx).
		Find(&orderItems).
		Error
	return
}
