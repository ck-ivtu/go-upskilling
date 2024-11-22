package su4

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func GuessingGame() {
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(101)

	fmt.Println("Guess the number between 0 and 100.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your guess: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred, try again.")
			continue
		}

		input = strings.TrimSpace(input)

		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input, enter an integer between 0 and 100.")
			continue
		}

		if guess < 0 || guess > 100 {
			fmt.Println("Please enter a number between 0 and 100.")
			continue
		}

		if guess < secretNumber {
			fmt.Println("Too low.")
		} else if guess > secretNumber {
			fmt.Println("Too high.")
		} else {
			fmt.Println("You guessed.")
			break
		}
	}
}
