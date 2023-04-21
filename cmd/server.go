package main

import (
	"gouser/config"
	"gouser/internal/server"
	"gouser/internal/server/handler"
	"gouser/utils/initialize"

	"go.uber.org/fx"
)

func serverRun() {
	app := fx.New(
		fx.Provide(
			// postgresql
			initialize.NewUserDB,
		),
		config.Module,
		initialize.Module,
		server.Module,
		handler.Module,
	)

	// Run app forever
	app.Run()
}
