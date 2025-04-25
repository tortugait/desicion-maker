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

	"github.com/tortugait/desicion-maker/internal/log"
)

func main() {
	// init main App context
	appCtx, cancelFunc := context.WithCancel(context.Background())

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
