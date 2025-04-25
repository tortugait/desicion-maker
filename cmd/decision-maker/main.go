package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tortugait/desicion-maker/internal/config"
	"github.com/tortugait/desicion-maker/internal/log"
	"github.com/tortugait/desicion-maker/internal/transport/http"
	httpHandler "github.com/tortugait/desicion-maker/internal/transport/http/handler"
)

func main() {
	if wd := getWorkDirArg(); wd != nil {
		if err := os.Chdir(*wd); err != nil {
			panic(err)
		}
	}

	// load App config
	conf, err := config.Load[config.App]()
	if err != nil {
		panic(err)
	}

	// init main App context
	appCtx, cancelFunc := context.WithCancel(context.Background())

	question := httpHandler.NewQuestion()
	httpSysHandler := httpHandler.NewSystem()
	handlers := http.Handlers{
		Status: httpSysHandler.GetStatus,
		Ask:    question.Ask,
	}
	httpSrv, err := initHTTPServer(conf.HTTPServer, handlers)
	if err != nil {
		log.Logger.Fatalw("init HTTP server", "error", err)
	}

	go func() {
		err := httpSrv.Run(appCtx)
		if err != nil {
			log.Logger.Fatalw("run HTTP server", "error", err)
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-osSignals
	log.Logger.Infof("received signal: %s. Canceling background jobs and exiting...", s)
	cancelFunc()
}

func initHTTPServer(conf config.HTTPServer, handlers http.Handlers) (*http.Server, error) {
	httpServer, err := http.NewServer(http.Config{
		DocsBase:        conf.DocsBase,
		HandlerTimeout:  conf.HandlerTimeout,
		ShutdownTimeout: conf.ShutdownTimeout,
		Addr:            conf.Addr,
	}, handlers)
	if err != nil {
		return nil, fmt.Errorf("create http server: %w", err)
	}

	return httpServer, nil
}

func getWorkDirArg() *string {
	w := flag.String("w", "", "Program working directory")
	flag.Parse()

	if *w == "" {
		return nil
	}
	return w
}
