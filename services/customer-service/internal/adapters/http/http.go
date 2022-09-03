package httpclient

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/fx"
)

type Client interface {
	Log(isEnable bool) *resty.Client
}

type ClientImpl struct {
	client *resty.Client
}

func NewClient() Client {
	client := resty.New()

	return &ClientImpl{
		client: client,
	}
}

func (c *ClientImpl) Log(isEnable bool) *resty.Client {
	return c.client.SetDebug(isEnable)
}

var Module = fx.Options(
	fx.Provide(NewClient),
	fx.Provide(NewOrderHttpClient),
)
