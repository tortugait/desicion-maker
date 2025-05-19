package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tortugait/desicion-maker/internal/service"
)

type questionRequest struct {
	Question string `json:"question"`
}

type questionResponse struct {
	Data any `json:"data"`
}

type question struct{}

func NewQuestion() question {
	return question{}
}

func (h question) Ask(eCtx echo.Context) error {
	var req questionRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	answer := "No"
	if service.ShouldDoIt() {
		answer = "Yes"
	}

	res := questionResponse{
		Data: map[string]string{
			"answer": answer,
		},
	}
	return eCtx.JSON(http.StatusOK, res)
}
