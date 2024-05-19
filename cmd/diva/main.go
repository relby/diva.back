package main

import (
	"context"

	"github.com/relby/diva.back/internal/app"
)

func main() {
	ctx := context.Background()
	app, err := app.NewApp(ctx)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
