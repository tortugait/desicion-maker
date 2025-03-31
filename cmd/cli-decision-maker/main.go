package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Decision-Maker CLI App")
	fmt.Println("-----------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Ask your question (or type 'exit' to quit): ")
		if !scanner.Scan() {
			fmt.Println("Error when reading input. Completion of work.")
			break
		}

		question := strings.TrimSpace(scanner.Text())

		// Exit the program
		if strings.EqualFold(question, "exit") {
			fmt.Println("Thank you for using the app. Goodbye!")
			break
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

func shouldDoIt() bool {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src) //nolint:gosec

	return rng.Intn(2) == 1
}
