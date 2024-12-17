package main

import (
	"fmt"
	"github.com/meli-fresh-products-api-backend-t1/internal/application"
)

func main() {
	// env
	// ...

	// app
	// - config
	cfg := &application.ConfigServerChi{
		ServerAddress: ":8080",
	}
	app := application.NewServerChi(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
