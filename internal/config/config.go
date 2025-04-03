package config

import "time"

//nolint:lll
type (
	App struct {
		HTTPServer `env-prefix:"HTTP_SERVER_"`
	}

	HTTPServer struct {
		DocsBase        string        `env:"DOCS_BASE" env-default:"http://localhost:8080" env-description:"Base docs path used to serve openAPI"`
		Addr            string        `env:"ADDR" env-default:":8080" env-description:"HTTP address to serve API requests"`
		HandlerTimeout  time.Duration `env:"HANDLER_TIMEOUT" env-default:"10s" env-description:"Timeout to handle request, zero means no timeout"`
		CheckTimeout    time.Duration `env:"CHECK_TIMEOUT" env-default:"10s" env-description:"Timeout to perform healthcheck, zero means no timeout"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"5s" env-description:"Timeout to gracefully shutdown API server"`
	}
)
