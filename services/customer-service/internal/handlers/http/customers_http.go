package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/usecases"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/pkg"
)

type CustomerHttpHandler interface {
	GetDetailProfile(c echo.Context) error
}

type CustomerHttpHandlerImpl struct {
	customerUsecase usecases.CustomerUsecase
}

func NewCustomerHttpHandler(customerUsecase usecases.CustomerUsecase) CustomerHttpHandler {
	return &CustomerHttpHandlerImpl{
		customerUsecase: customerUsecase,
	}
}

func (h *CustomerHttpHandlerImpl) GetDetailProfile(c echo.Context) error {
	var (
		err      error
		params   *models.DetailProfileRequest = new(models.DetailProfileRequest)
		response pkg.Response
		data     models.DetailProfileResponse
	)

	if err = c.Bind(params); err != nil {
		response = pkg.NewResponse(data, pkg.NewHTTPError(http.StatusBadRequest, err))
		return c.JSON(response.Code, response)
	}

	ctx := pkg.NewRequestFromHeader(c.Request().Header, c.Request().Context())

	ctx, span := pkg.NewSpan(ctx, "CustomerHttpHandlerImpl.GetDetailProfile", nil)
	defer span.End()

	data, err = h.customerUsecase.GetDetailProfile(ctx, params)
	response = pkg.NewResponse(data, err)

	return c.JSON(response.Code, response)
}
