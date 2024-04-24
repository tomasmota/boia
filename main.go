package main

import (
	"fmt"
)

const (
	PeerPort   = "3838"
	ClientPort = "3839"
)

func main() {
	log := NewLog()
	sm := NewStateMachine()

	log.Append(LogEntry{
		Term: 233123,
		Command: &SetCommand{
			Key:   "abc",
			Value: 32,
		},
		Committed: false,
	})

	log.Entries[0].Command.Apply(sm)

	fmt.Println(sm.state["abc"])
}
