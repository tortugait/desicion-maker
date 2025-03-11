package main

import (
	"fmt"
)

func main() {
	fmt.Println("Decision-Maker CLI App")
	fmt.Println("-----------------------")

	for {
		fmt.Print("Ask your question (or type 'exit' to quit): ")
		// todo: scan user input

		// todo: handle 'exit' command

		// todo: handle empty question ('')

		if shouldDoIt() {
			// todo: print the positive answer
		}

		// todo: print the negative answer
	}
}

func shouldDoIt() bool {
	// todo: implement, return true or false randomly
	return false
}
