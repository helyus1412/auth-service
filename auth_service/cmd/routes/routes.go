package routes

import (
	"net/http"

	"github.com/helyus1412/auth-service/domain/auth"
	"github.com/helyus1412/auth-service/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

func InitRoutes(e *echo.Echo, db *sqlx.DB, tc trace.Tracer, logger *logger.Logger) {
	e.GET("/health-check", func(e echo.Context) error {
		return e.String(http.StatusOK, "auth service is running properly")
	})

	authRepository := auth.NewRepository(db, "")
	authUsecase := auth.NewUsecase(authRepository)
	authHandler := auth.NewHandler(authUsecase, tc)

	e.POST("/register", authHandler.Register)
}
