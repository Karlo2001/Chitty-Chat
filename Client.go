package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var clock []int
var input chan int = make(chan int, 1)
var message string
var user string
var id int
var scanner = bufio.NewScanner(os.Stdin)

func main() {
	fmt.Println("Please enter a username")
	scanner.Scan()
	user = scanner.Text()

	//Connect to server
	//Get id from server
	fmt.Scanln(&id)
	clock = make([]int, id+1)
	fmt.Println("You are now connected and can type messages")
	for {
		scanner.Scan()
		message = scanner.Text()
		clock[id]++
		newMessage := Message{sender: user, message: message, clock: clock}

		publish(newMessage)
		readMessage(newMessage)
		fmt.Println(clock)
	}
}

func publish(message Message) {
	if validateMessage(message) {
		//Send message
	} else {
		fmt.Println("The string provided is too long. The allowed amount is 128 yours was " + strconv.Itoa(len(message.message)))
	}
}

func readMessage(message Message) {
	updateClock(message.clock)
	logMessage(message)
}

func updateClock(otherClock []int) {
	for i, s := range otherClock {
		if i >= len(clock) {
			clock = append(clock, s)
		} else {
			clock[i] = max(clock[i], s)
		}
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func lamportTime(clock []int) int {
	var l int
	for _, s := range clock {
		l += s
	}
	return l
}
