package main

import (
	"log"

	"github.com/helyus1412/auth-service/pkg/databases"
	"github.com/pressly/goose/v3"
)

func main() {
	db, err := databases.InitPostgre()
	if err != nil {
		log.Fatalf("failed init postgre: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(db.DB, "./pkg/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrated successfully")
}
