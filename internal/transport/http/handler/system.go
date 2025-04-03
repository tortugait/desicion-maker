package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	system struct{}
)

func NewSystem() system {
	return system{}
}

func (h system) GetStatus(eCtx echo.Context) error {
	response := map[string]string{"status": "OK"}
	return eCtx.JSON(http.StatusOK, response)
}
