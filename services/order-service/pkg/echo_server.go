package pkg

import "github.com/labstack/echo/v4"

// EchoServer echo server
type EchoServer struct {
	Echo *echo.Echo
}

// NewEchoServer create echo server
func NewEchoServer() *EchoServer {
	e := echo.New()
	return &EchoServer{Echo: e}
}
