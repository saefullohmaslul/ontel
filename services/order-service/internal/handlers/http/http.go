package httphandler

import (
	"github.com/saefullohmaslul/distributed-tracing/order-service/pkg"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewOrderHttpHandler),
	fx.Provide(OrderRoute),
)

type Route interface {
	Setup()
}

type OrderRouteImpl struct {
	echo  *pkg.EchoServer
	order OrderHttpHandler
}

func OrderRoute(echo *pkg.EchoServer, order OrderHttpHandler) Route {
	return &OrderRouteImpl{
		echo:  echo,
		order: order,
	}
}

func (r *OrderRouteImpl) Setup() {
	order := r.echo.Echo.Group("/v1")
	{
		order.GET("/orders/:customer_id", r.order.GetOrdersDetail)
	}
}
