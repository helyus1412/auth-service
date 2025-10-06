package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/helyus1412/auth-service/cmd/routes"
	"github.com/helyus1412/auth-service/config"
	"github.com/helyus1412/auth-service/pkg/databases"
	"github.com/labstack/echo/v4"
	"github.com/pressly/goose/v3"
)

func main() {
	e := echo.New()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := databases.InitPostgre()
	if err != nil {
		log.Fatalf("failed init postgre: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}
	if err := goose.Up(db.DB, "./pkg/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("Database migrated successfully")

	routes.InitRoutes(e, db)

	go func() {
		port := fmt.Sprintf(":%d", config.GlobalEnv.HTTPPort)
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed start server: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelShutdown()

	// Shutdown the HTTP server.
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("failed shutdown server: %v", err)
	}
}
