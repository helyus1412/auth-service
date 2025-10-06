package databases

import (
	"fmt"
	"time"

	"github.com/helyus1412/auth-service/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitPostgre() (*sqlx.DB, error) {
	dbConfig := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.GlobalEnv.PostgreHost, config.GlobalEnv.PostgreUser, config.GlobalEnv.PostgrePassword,
		config.GlobalEnv.PostgreDBName, config.GlobalEnv.PostgrePort, config.GlobalEnv.PostgreSSLMode)
	conn, err := sqlx.Connect("postgres", dbConfig)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(config.GlobalEnv.PostgreMaxOpenCon)
	conn.SetMaxIdleConns(config.GlobalEnv.PostgreMaxIddleCon)
	conn.SetConnMaxLifetime(time.Minute * time.Duration(config.GlobalEnv.PostgreMaxLifeTime))

	return conn, nil
}
