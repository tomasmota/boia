package main

type Log struct {
	Entries []LogEntry

	// Highest index we know to be commited
	// This is sent in AppendEntries
	HighestCommited int64
}

type LogEntry struct {
	Term     int64
	Command  StateMachineCommand
	Commited bool
}

type StateMachineCommand struct {
	Key   string
	Value int
}

// func (l *Log) AddEntry(entry LogEntry) {
// 	l.Entries = append(l.Entries, entry)
// }
