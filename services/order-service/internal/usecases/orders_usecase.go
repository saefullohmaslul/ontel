package usecases

import (
	"context"
	"net/http"
	"time"

	"github.com/saefullohmaslul/distributed-tracing/order-service/internal/adapters/postgres"
	"github.com/saefullohmaslul/distributed-tracing/order-service/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/order-service/pkg"
)

type OrderUsecase interface {
	GetOrders(ctx context.Context, params *models.OrdersRequest) (data []models.Order, err error)
	GetOrdersDetail(ctx context.Context, params *models.OrdersRequest) (data []models.OrdersDetailResponse, err error)
}

type OrderUsecaseImpl struct {
	OrderPostgres postgres.OrderPostgres
}

func NewOrderUsecase(OrderPostgres postgres.OrderPostgres) OrderUsecase {
	return &OrderUsecaseImpl{
		OrderPostgres: OrderPostgres,
	}
}

func (u *OrderUsecaseImpl) GetOrders(ctx context.Context, params *models.OrdersRequest) (data []models.Order, err error) {
	ctx, span := pkg.NewSpan(ctx, "OrderUsecaseImpl.GetOrders", nil)
	defer span.End()

	data, err = u.OrderPostgres.GetOrders(ctx, &models.OrdersRequest{
		CustomerID: params.CustomerID,
	})

	if err != nil {
		err = pkg.NewHTTPError(http.StatusInternalServerError, err)
	}

	return
}

func (u *OrderUsecaseImpl) GetOrderItems(ctx context.Context, params *models.OrderItemsRequest) (data []models.OrderItem, err error) {
	ctx, span := pkg.NewSpan(ctx, "OrderUsecaseImpl.GetOrderItems", nil)
	defer span.End()

	data, err = u.OrderPostgres.GetOrderItems(ctx, &models.OrderItemsRequest{
		OrderID: params.OrderID,
	})

	if err != nil {
		err = pkg.NewHTTPError(http.StatusInternalServerError, err)
	}

	time.Sleep(5 * time.Second)

	return
}

func (u *OrderUsecaseImpl) GetOrdersDetail(ctx context.Context, params *models.OrdersRequest) (data []models.OrdersDetailResponse, err error) {
	ctx, span := pkg.NewSpan(ctx, "OrderUsecaseImpl.GetOrdersDetail", nil)
	defer span.End()

	orders, err := u.GetOrders(ctx, params)

	if err != nil {
		return
	}

	data = []models.OrdersDetailResponse{}

	for _, order := range orders {
		items, err := u.GetOrderItems(ctx, &models.OrderItemsRequest{
			OrderID: order.OrderID,
		})

		if err != nil {
			return data, err
		}

		data = append(data, models.OrdersDetailResponse{
			Order: order,
			Items: items,
		})
	}

	return
}
