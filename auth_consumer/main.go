package main

import (
	"fmt"
	"log"
	"os"

	authPkg "github.com/helyus1412/auth-service/domain/auth"
	authDto "github.com/helyus1412/auth-service/dto"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := InitPostgre()
	if err != nil {
		log.Fatalf("failed init postgre: %v", err)
	}

	//from auth_service
	authRepo := authPkg.NewRepository(conn, "")
	authUsecase := authPkg.NewUsecase(authRepo)

	err = authUsecase.Register(&authDto.RegisterRequest{
		Email:    "helmi",
		Password: "helmi",
	})

	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Print("success register user")
	}
}

func InitPostgre() (*sqlx.DB, error) {
	var ok bool

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	postgreHost, ok := os.LookupEnv("POSTGRE_HOST")
	if !ok {
		panic("missing POSTGRE_HOST environment")
	}

	postgreUser, ok := os.LookupEnv("POSTGRE_USER")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	postgrePassword, ok := os.LookupEnv("POSTGRE_PASSWORD")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	postgreDBName, ok := os.LookupEnv("POSTGRE_DB_NAME")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	postgrePort, ok := os.LookupEnv("POSTGRE_PORT")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	postgreSSLMode, ok := os.LookupEnv("POSTGRE_SSL_MODE")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	dbConfig := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		postgreHost, postgreUser, postgrePassword,
		postgreDBName, postgrePort, postgreSSLMode)

	conn, err := sqlx.Connect("postgres", dbConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
