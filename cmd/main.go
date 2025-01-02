package main

import (
	"fmt"
	"github.com/go-sql-driver/mysql"

	"github.com/meli-fresh-products-api-backend-t1/internal/application"
)

func main() {
	mysqlCfg := mysql.Config{
		User:      "melisprint_user",
		Passwd:    "melisprint_pass",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "melisprint",
		ParseTime: true,
	}
	cfg := &application.ConfigServerChi{
		ServerAddress: ":8080",
		Dsn:           mysqlCfg.FormatDSN(),
	}
	server := application.NewServerChi(cfg)
	fmt.Println("Server running on port 8080...")

	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
