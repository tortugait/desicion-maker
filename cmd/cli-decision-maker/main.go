package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tortugait/decision-maker/internal/log"
	"github.com/tortugait/decision-maker/internal/service"
)

func main() {
	// init main App context
	appCtx, cancelFunc := context.WithCancel(context.Background())
	decisionSrv := service.NewDecisionSrv(service.NewDeepSeekSrv(service.DeepSeekConf{}))

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
					log.Logger.Fatal("error when reading input, completion of work")
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

				if decisionSrv.ShouldDoIt() {
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
