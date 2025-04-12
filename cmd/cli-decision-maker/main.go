package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tortugait/desicion-maker/internal/config"
	"github.com/tortugait/desicion-maker/internal/log"
	"github.com/tortugait/desicion-maker/internal/transport/http"
	httpHandler "github.com/tortugait/desicion-maker/internal/transport/http/handler"
)

func main() {
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

	go func() {
		fmt.Println("Decision-Maker CLI App")
		fmt.Println("-----------------------")
		scanner := bufio.NewScanner(os.Stdin)

		for {
			select {
			case <-appCtx.Done():
				return
			default:
				fmt.Print("Ask your question (or type 'exit' to quit): \n")
				if !scanner.Scan() {
					log.Logger.Fatalw("error when reading input, completion of work", "error", err)
					os.Exit(1)
				}

				question := strings.TrimSpace(scanner.Text())

				// Exit the program
				if strings.EqualFold(question, "exit") {
					fmt.Println("Thank you for using the app. Goodbye!")
					os.Exit(0)
				}

				if question == "" {
					continue
				}

				if shouldDoIt() {
					fmt.Println("Yes")
				} else {
					fmt.Println("No")
				}
			}
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-osSignals
	log.Logger.Infof("received signal: %s. Canceling background jobs and exiting...", s)
	cancelFunc()
}

func shouldDoIt() bool {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src) //nolint:gosec

	return rng.Intn(2) == 1 //nolint:gomnd
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
