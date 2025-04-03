package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestSystem_GetStatus(t *testing.T) {
	t.Parallel()

	e := echo.New()

	testCases := []struct {
		name           string
		expectErr      error
		expectedStatus int
	}{
		{
			name:           "status OK",
			expectErr:      nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/status", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			eCtx := e.NewContext(req, rec)

			err := NewSystem().GetStatus(eCtx)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
