package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var GlobalEnv Env

type Env struct {
	HTTPPort             uint16
	PostgreHost          string
	PostgreUser          string
	PostgrePassword      string
	PostgreDBName        string
	PostgrePort          string
	PostgreSSLMode       string
	PostgreMaxOpenCon    int
	PostgreMaxIddleCon   int
	PostgreMaxLifeTime   int
	PostgreDefaultSchema string
	AppName              string
	AppVersion           string
	IsProduction         bool
	OtelTraceHost        string
	OtelTracePort        int
	OtelMetricsHost      string
	OtelMetricsPort      int
}

func init() {
	var ok bool

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	if port, err := strconv.Atoi(os.Getenv("HTTP_PORT")); err != nil {
		panic("missing HTTP_PORT environment")
	} else {
		GlobalEnv.HTTPPort = uint16(port)
	}

	GlobalEnv.PostgreHost, ok = os.LookupEnv("POSTGRE_HOST")
	if !ok {
		panic("missing POSTGRE_HOST environment")
	}

	GlobalEnv.PostgreUser, ok = os.LookupEnv("POSTGRE_USER")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	GlobalEnv.PostgrePassword, ok = os.LookupEnv("POSTGRE_PASSWORD")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	GlobalEnv.PostgreDBName, ok = os.LookupEnv("POSTGRE_DB_NAME")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	GlobalEnv.PostgrePort, ok = os.LookupEnv("POSTGRE_PORT")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	GlobalEnv.PostgreSSLMode, ok = os.LookupEnv("POSTGRE_SSL_MODE")
	if !ok {
		panic("missing POSTGRE_USER environment")
	}

	if postgreMaxOpenCon, err := strconv.Atoi(os.Getenv("POSTGRE_MAX_OPEN_CON")); err != nil {
		panic("missing POSTGRE_MAX_CON environment")
	} else {
		GlobalEnv.PostgreMaxOpenCon = int(postgreMaxOpenCon)
	}

	if postgreMaxIddleCon, err := strconv.Atoi(os.Getenv("POSTGRE_MAX_IDDLE_CON")); err != nil {
		panic("missing POSTGRE_MAX_CON environment")
	} else {
		GlobalEnv.PostgreMaxIddleCon = int(postgreMaxIddleCon)
	}

	if postgreMaxLifeTime, err := strconv.Atoi(os.Getenv("POSTGRE_MAX_LIFE_TIME")); err != nil {
		panic("missing POSTGRE_MAX_CON environment")
	} else {
		GlobalEnv.PostgreMaxLifeTime = int(postgreMaxLifeTime)
	}

	GlobalEnv.PostgreDefaultSchema, ok = os.LookupEnv("POSTGRE_DEFAULT_SCHEMA")
	if !ok {
		GlobalEnv.PostgreDefaultSchema = ""
	}

	GlobalEnv.AppName, ok = os.LookupEnv("APP_NAME")
	if !ok {
		log.Panicln("config.init() missing APP_NAME environment")
	}

	GlobalEnv.AppVersion, ok = os.LookupEnv("APP_VERSION")
	if !ok {
		log.Panicln("config.init() missing APP_VERSION environment")
	}

	isProductionStr, ok := os.LookupEnv("IS_PRODUCTION")
	if !ok {
		panic("missing IS_PRODUCTION environment")
	}

	isProduction, err := strconv.ParseBool(isProductionStr)
	if err != nil {
		panic("invalid value for IS_PRODUCTION, must be true or false")
	}

	GlobalEnv.IsProduction = isProduction

	GlobalEnv.OtelTraceHost, ok = os.LookupEnv("OTEL_TRACE_HOST")
	if !ok {
		log.Panicln("config.init() missing OTEL_TRACE_HOST environment")
	}

	if port, err := strconv.Atoi(os.Getenv("OTEL_TRACE_PORT")); err != nil {
		panic("missing OTEL_TRACE_PORT environment")
	} else {
		GlobalEnv.OtelTracePort = port
	}

	GlobalEnv.OtelMetricsHost, ok = os.LookupEnv("OTEL_METRICS_HOST")
	if !ok {
		log.Panicln("config.init() missing OTEL_METRICS_HOST environment")
	}

	if port, err := strconv.Atoi(os.Getenv("OTEL_METRICS_PORT")); err != nil {
		panic("missing OTEL_METRICS_PORT environment")
	} else {
		GlobalEnv.OtelMetricsPort = port
	}
}
