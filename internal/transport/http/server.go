package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/tortugait/desicion-maker/internal/log"
)

const (
	docsPath         = "./api/http"
	docsTemplatePath = docsPath + "/index_tmpl.html"
	docsOutPath      = docsPath + "/index.html"
)

type (
	Config struct {
		DocsBase        string
		HandlerTimeout  time.Duration
		ShutdownTimeout time.Duration
		Addr            string
		AuthToken       string
	}

	Handlers struct {
		Status              echo.HandlerFunc
		AddActivePair       echo.HandlerFunc
		UpdateActivePair    echo.HandlerFunc
		ListActivePair      echo.HandlerFunc
		AddCandidatePair    echo.HandlerFunc
		UpdateCandidatePair echo.HandlerFunc
		ListCandidatePair   echo.HandlerFunc
		ListExchange        echo.HandlerFunc
	}

	Server struct {
		echo       *echo.Echo
		checkRoute string
		config     Config
	}
)

func NewServer(config Config, handlers Handlers) (*Server, error) {
	if err := prepareDocs(config.DocsBase); err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(config.Addr)
	if err != nil {
		return nil, fmt.Errorf("split addr: %w", err)
	}
	// if server listens on all interfaces, run checks against localhost
	if host == "" || host == "0.0.0.0" {
		host = "localhost"
	}

	s := &Server{
		echo:       echo.New(),
		config:     config,
		checkRoute: fmt.Sprintf("http://%s%s", net.JoinHostPort(host, port), checkRoutePath),
	}

	s.echo.Server.ReadTimeout = s.config.HandlerTimeout
	s.echo.Server.WriteTimeout = s.config.HandlerTimeout
	s.echo.HideBanner = true
	s.echo.HidePort = true
	s.echo.Validator = &customValidator{validator: validator.New()}

	InitRoutes(s.echo, handlers)
	return s, nil
}

// Run starts HTTP server. This is a blocking call.
// To stop serving, cancel the passed context.
func (s *Server) Run(ctx context.Context) error {
	e := make(chan error, 1)
	go func() {
		e <- s.echo.Start(s.config.Addr)
	}()

	log.Logger.Infow(
		"HTTP server is running",
		"addr", s.config.Addr,
	)

	select {
	case <-ctx.Done():
		ctxWithTimeout, cancelCtxWithTimeout := context.WithTimeout(ctx, s.config.ShutdownTimeout)
		defer cancelCtxWithTimeout()
		if err := s.echo.Shutdown(ctxWithTimeout); err != nil {
			return err
		}
		return nil
	case err := <-e:
		return err
	}
}

func (s *Server) Check(ctx context.Context) error {
	req, err := http.NewRequest(http.MethodGet, s.checkRoute, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("do check request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return errors.New("http server is not available at the moment") //nolint:err113
	}

	return nil
}

func prepareDocs(basePath string) error {
	t, err := template.New(path.Base(docsTemplatePath)).ParseFiles(docsTemplatePath)
	if err != nil {
		return fmt.Errorf("open docs template: %w", err)
	}

	out, err := os.Create(docsOutPath)
	if err != nil {
		return fmt.Errorf("open docs output for writing: %w", err)
	}
	defer out.Close() //nolint:errcheck

	if err := t.Execute(out, map[string]string{
		"base_path": basePath + v1DocsRoutePath,
	}); err != nil {
		return fmt.Errorf("execute docs template: %w", err)
	}

	return nil
}
