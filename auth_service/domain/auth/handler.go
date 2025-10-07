package auth

import (
	"net/http"

	"github.com/helyus1412/auth-service/dto"
	httpError "github.com/helyus1412/auth-service/pkg/httpError"
	"github.com/helyus1412/auth-service/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

type Handler interface {
	Register(c echo.Context) error
}

type handler struct {
	usecase Usecase
	tc      trace.Tracer
}

func NewHandler(usecase Usecase, tc trace.Tracer) Handler {
	return &handler{usecase, tc}
}

func (h *handler) Register(c echo.Context) error {
	ctx, span := h.tc.Start(c.Request().Context(), "handler.Register")
	defer span.End()

	var payload dto.RegisterRequest

	if err := c.Bind(&payload); err != nil {
		respErr := httpError.NewBadRequest()
		return utils.ResponseError(respErr, c)
	}

	result := h.usecase.Register(ctx, &payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Success Register", http.StatusCreated, c)
}
