// @title Meli Fresh Products API
// @version 1.0
// @description API for managing fresh products and orders
// @contact.name Bootcampers GO - W5
// @host localhost:8080
// @BasePath /

package main

import (
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"

	"github.com/meli-fresh-products-api-backend-t1/internal/application"
)

func main() {
	dbUri := os.Getenv("MYSQL_SPRINT_URI")
	if dbUri == "" {
		dbUri = "localhost:3306"
	}
	mysqlCfg := mysql.Config{
		User:      "root",
		Passwd:    "meli_pass",
		Net:       "tcp",
		Addr:      dbUri,
		DBName:    "melifresh",
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
