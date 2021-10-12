package main

import (
	"log"
	"strconv"
)

type Message struct {
	sender  string
	message string
	clock   []int
	id      int
}

func logMessage(message Message) {
	log.Println(message.sender + ": " + message.message + " Lamport Time: " + strconv.Itoa(lamportTime(message.clock)))
}

func validateMessage(message Message) bool {
	if len(message.message) > 128 {
		return false

	}
	return true
}
