package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/saefullohmaslul/distributed-tracing/internal/adapters/postgres"
	"github.com/saefullohmaslul/distributed-tracing/pkg"
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
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	echoServer *pkg.EchoServer,
	postgresDatabase *postgres.Database,
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
