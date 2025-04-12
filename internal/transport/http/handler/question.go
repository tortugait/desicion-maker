package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	answer := "No"
	if rng.Intn(2) == 1 {
		answer = "Yes"
	}

	res := questionResponse{
		Data: map[string]string{
			"answer": answer,
		},
	}
	return eCtx.JSON(http.StatusOK, res)
}
