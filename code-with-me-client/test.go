package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter a number: ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1] // Remove the trailing newline character

		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			return
		}

		result := number * 2

		fmt.Printf("The result is: %d\n", result)
	}

}
