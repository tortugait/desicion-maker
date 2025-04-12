package handler

import (
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"time"
)

type QuestionRequest struct {
	Question string `json:"question"`
}

type question struct{}

func NewQuestion() question {
	return question{}
}

func (h question) Ask(eCtx echo.Context) error {
	var req QuestionRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	answer := "No"
	if rng.Intn(2) == 1 {
		answer = "Yes"
	}

	return eCtx.JSON(http.StatusOK, map[string]string{
		"answer": answer,
	})
}
