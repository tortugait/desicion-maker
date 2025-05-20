package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tortugait/desicion-maker/internal/log"
)

const (
	checkRoutePath = "/check"

	v1APIRoutePrefix = "/api/v1"
	v1DocsRoutePath  = "/api/v1/docs"
)

func InitRoutes(e *echo.Echo, handlers Handlers) {
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	rootGroup := e.Group("")
	rootGroup.GET(checkRoutePath, func(eCtx echo.Context) error { return nil })
	rootGroup.Static(v1DocsRoutePath, docsPath)

	v1 := rootGroup.Group(
		v1APIRoutePrefix,
		newLoggingMiddleware(log.Logger),
	)

	v1.GET("/status", handlers.Status)
	v1.POST("/ask", handlers.Ask)
}
