package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saefullohmaslul/distributed-tracing/order-service/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/order-service/internal/usecases"
	"github.com/saefullohmaslul/distributed-tracing/order-service/pkg"
)

type OrderHttpHandler interface {
	GetOrdersDetail(c echo.Context) error
}

type OrderHttpHandlerImpl struct {
	OrderUsecase usecases.OrderUsecase
}

func NewOrderHttpHandler(OrderUsecase usecases.OrderUsecase) OrderHttpHandler {
	return &OrderHttpHandlerImpl{
		OrderUsecase: OrderUsecase,
	}
}

func (h *OrderHttpHandlerImpl) GetOrdersDetail(c echo.Context) error {
	var (
		err      error
		params   *models.OrdersDetailRequest = new(models.OrdersDetailRequest)
		response pkg.Response
		data     []models.OrdersDetailResponse
	)

	if err = c.Bind(params); err != nil {
		response = pkg.NewResponse(data, pkg.NewHTTPError(http.StatusBadRequest, err))
		return c.JSON(response.Code, response)
	}

	if err = (&echo.DefaultBinder{}).BindPathParams(c, params); err != nil {
		response = pkg.NewResponse(data, pkg.NewHTTPError(http.StatusBadRequest, err))
		return c.JSON(response.Code, response)
	}

	data, err = h.OrderUsecase.GetOrdersDetail(c.Request().Context(), &models.OrdersRequest{
		CustomerID: params.CustomerID,
	})

	response = pkg.NewResponse(data, err)

	return c.JSON(response.Code, response)
}
