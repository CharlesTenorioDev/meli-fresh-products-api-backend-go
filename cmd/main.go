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
	fmt.Println("Server running on port 8080...")

	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
