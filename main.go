package main

import (
	"main/individuals"
	"main/migrations"
	"main/professions"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	professions.RegisterRoutes(app)
	individuals.RegisterRoutes(app)

	app.Run()
}
