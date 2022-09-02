package httphandler

import (
	"github.com/saefullohmaslul/distributed-tracing/customer-service/pkg"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewCustomerHttpHandler),
	fx.Provide(NewCustomerRoute),
)

type Route interface {
	Setup()
}

type CustomerRouteImpl struct {
	echo     *pkg.EchoServer
	customer CustomerHttpHandler
}

func NewCustomerRoute(echo *pkg.EchoServer, customer CustomerHttpHandler) Route {
	return &CustomerRouteImpl{
		echo:     echo,
		customer: customer,
	}
}

func (r *CustomerRouteImpl) Setup() {
	customer := r.echo.Echo.Group("/v1")
	{
		customer.GET("/profile/:customer_id", r.customer.GetDetailProfile)
	}
}
