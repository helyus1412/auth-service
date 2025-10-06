package auth

import (
	"net/http"

	"github.com/helyus1412/auth-service/dto"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Register(c echo.Context) error
}

type handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return &handler{usecase}
}

func (h *handler) Register(c echo.Context) error {
	var payload dto.RegisterRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid json format",
		})
	}

	err := h.usecase.Register(&payload)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, "success register user")
}
