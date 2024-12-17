package main

import (
	"fmt"

	"github.com/meli-fresh-products-api-backend-t1/internal/application"
)

func main() {
	cfg := &application.ConfigServerChi{
		ServerAddress: ":8080",
	}
	server := application.NewServerChi(cfg)

	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
