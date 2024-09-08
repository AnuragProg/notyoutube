package main

import (
	"fmt"
	"os"
	"time"
)

func Clear() {
	fmt.Print("\033[2J")
}

func UpdateExistingFrameExample() {
	characters := []string{"a", "b", "c", "d", "e", "f"}
	current := 0
	Clear()
	fmt.Print("\033[H") // go to home position
	fmt.Print("HelloWorld")
	for {
		fmt.Print("\033[H") // go to home position
		fmt.Print(characters[current])
		current = (current + 1) % len(characters)
		time.Sleep(time.Second * 2)
	}
}

func ExampleOfGettingCursorPosition() {
	// maintain internal state to keep track of the state so that 
	// time is not wasted and complexity is prevented just for knowing where the 
	// cursor is when we are the ones moving it
}

func main() {
	ExampleOfGettingCursorPosition()
}
