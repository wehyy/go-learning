package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	maxNum := 100
	var yourGuess int
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(maxNum)

	fmt.Println("Please input your guess")
	for {
		n, err := fmt.Scanf("%d", &yourGuess)

		if err != nil || n != 1 {
			fmt.Println("Invalid input. Please enter an integer value")
			continue
		}
		fmt.Println("Your guess is ", yourGuess)
		if yourGuess > secretNumber {
			fmt.Println("Your guess is bigger than the secret number. Please try again")
		} else if yourGuess < secretNumber {
			fmt.Println("Your guess is smaller than the secret number. Please try again")
		} else {
			fmt.Println("Correct, you legend!")
			break
		}
	}
}
