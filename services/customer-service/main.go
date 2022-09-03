package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	httpclient "github.com/saefullohmaslul/distributed-tracing/customer-service/internal/adapters/http"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/adapters/postgres"
	httphandler "github.com/saefullohmaslul/distributed-tracing/customer-service/internal/handlers/http"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/internal/usecases"
	"github.com/saefullohmaslul/distributed-tracing/customer-service/pkg"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load()
	fx.New(Module).Run()
}

// Module main module
var Module = fx.Options(
	pkg.Module,
	postgres.Module,
	httphandler.Module,
	usecases.Module,
	httpclient.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	echoServer *pkg.EchoServer,
	postgresDatabase *postgres.Database,
	routes httphandler.Route,
) {
	conn, _ := postgresDatabase.DB.DB()
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					port, found := os.LookupEnv("PORT")

					if !found {
						port = "7150"
					}

					routes.Setup()

					prv, err := pkg.NewProvider(pkg.TracerProviderConfig{
						JaegerEndpoint: os.Getenv("JAEGER_ENDPOINT"),
						ServiceName:    os.Getenv("SERVICE_NAME"),
						ServiceVersion: os.Getenv("SERVICE_VERSION"),
						Environment:    os.Getenv("ENVIRONMENT"),
						Disabled:       false,
					})

					if err != nil {
						return
					}

					defer prv.Close(ctx)

					echoServer.Echo.Logger.Fatal(
						echoServer.Echo.Start(fmt.Sprintf(":%s", port)),
					)
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return conn.Close()
			},
		},
	)
}
