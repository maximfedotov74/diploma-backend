package main

import (
	"github.com/maximfedotov74/fiber-psql/internal/app"
)

func main() {
	app := app.NewApplication()
	app.Start()
}
