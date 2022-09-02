package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/saefullohmaslul/distributed-tracing/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/internal/usecases"
	"github.com/saefullohmaslul/distributed-tracing/pkg"
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

	data, err = h.customerUsecase.GetDetailProfile(c.Request().Context(), params)
	response = pkg.NewResponse(data, err)

	return c.JSON(response.Code, response)
}
