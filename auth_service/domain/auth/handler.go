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
	Login(c echo.Context) error
	ListUser(c echo.Context) error
	Edit(c echo.Context) error
	Delete(c echo.Context) error
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
		respErr := httpError.NewBadRequest("")
		return utils.ResponseError(respErr, c)
	}

	result := h.usecase.Register(ctx, &payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Success Register", http.StatusCreated, c)
}

func (h *handler) Login(c echo.Context) error {
	ctx, span := h.tc.Start(c.Request().Context(), "handler.Login")
	defer span.End()

	var payload dto.LoginRequest

	if err := c.Bind(&payload); err != nil {
		respErr := httpError.NewBadRequest("")
		return utils.ResponseError(respErr, c)
	}

	result := h.usecase.Login(ctx, &payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Success Login", http.StatusOK, c)
}

func (h *handler) ListUser(c echo.Context) error {
	ctx, span := h.tc.Start(c.Request().Context(), "handler.ListUser")
	defer span.End()

	result := h.usecase.ListUser(ctx)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "List User", http.StatusOK, c)
}

func (h *handler) Edit(c echo.Context) error {
	ctx, span := h.tc.Start(c.Request().Context(), "handler.Edit")
	defer span.End()

	var payload dto.EditRequest

	if err := c.Bind(&payload); err != nil {
		respErr := httpError.NewBadRequest("")
		return utils.ResponseError(respErr, c)
	}

	result := h.usecase.Edit(ctx, &payload)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Success Edit", http.StatusCreated, c)
}

func (h *handler) Delete(c echo.Context) error {
	ctx, span := h.tc.Start(c.Request().Context(), "handler.Delete")
	defer span.End()

	var payload dto.DeleteRequest

	if err := c.Bind(&payload); err != nil {
		respErr := httpError.NewBadRequest("")
		return utils.ResponseError(respErr, c)
	}

	result := h.usecase.Delete(ctx, payload.ID)
	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Success Delete User", http.StatusOK, c)
}
