package main

import (
	"gouser/config"
	"gouser/internal/server"
	"gouser/internal/server/handler"
	"gouser/pkg/user"
	"gouser/utils/initialize"

	"go.uber.org/fx"
)

func serverRun() {
	app := fx.New(
		fx.Provide(
			// postgresql
			initialize.NewGoUserDB,
		),
		config.Module,
		initialize.Module,
		server.Module,
		handler.Module,
		user.Module,
	)

	// Run app forever
	app.Run()
}
