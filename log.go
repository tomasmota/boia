package main

import (
	"sync"
)

// LogEntryType describes various types of log entries.
type LogEntryType uint8

const (
	// LogCommand is applied to a user FSM.
	LogCommand LogEntryType = iota

	// LogNoop is used to assert leadership.
	LogNoop
)

type Log struct {
	mu      sync.RWMutex
	Entries []LogEntry

	// Highest index we know to be commited
	// This is sent in AppendEntries
	CommitIndex int
}

type LogEntry struct {
	// Index of the log entry
	Index uint64

	// Term when this entry was created
	Term uint64

	// Type of the entry
	Type LogEntryType

	// Whether the entry has been commited
	Committed bool

	// Command to be exectuted when entry is applied to the state machine
	Command Command
}

func NewLog() *Log {
	return &Log{
		Entries:     make([]LogEntry, 0, 1000),
		CommitIndex: -1,
	}
}

func (l *Log) Append(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Entries = append(l.Entries, entry)
}

func (l *Log) CommitEntry(idx int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Entries[idx].Committed = true
	l.CommitIndex = idx
}
