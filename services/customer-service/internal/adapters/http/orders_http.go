package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/models"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/pkg"
)

type OrderHttpClient interface {
	GetDetailOrders(ctx context.Context, params *models.OrdersDetailRequest) (data []models.OrdersDetailResponse, err error)
}

type OrderHttpClientImpl struct {
	client Client
}

func NewOrderHttpClient(client Client) OrderHttpClient {
	return &OrderHttpClientImpl{
		client: client,
	}
}

func (c *OrderHttpClientImpl) GetDetailOrders(ctx context.Context, params *models.OrdersDetailRequest) (data []models.OrdersDetailResponse, err error) {
	var response pkg.Response

	uri := fmt.Sprintf("%s/v1/orders/%d", os.Getenv("ORDER_API"), params.CustomerID)

	resp, err := c.client.
		Log(true).
		R().
		SetHeader("Content-Type", "application/json").
		SetResult(&response).
		SetError(&response).
		SetContext(ctx).
		Get(uri)

	if err != nil {
		return
	}

	if resp.IsError() {
		return
	}

	dataMarshal, err := json.Marshal(response.Data)
	if err != nil {
		return
	}

	err = json.Unmarshal(dataMarshal, &data)
	return
}
