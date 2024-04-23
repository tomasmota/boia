package main

type Command interface {
	Apply(*StateMachine) error
}

type StateMachine struct {
	state map[string]int
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		state: make(map[string]int),
	}
}

type SetCommand struct {
	Key   string
	Value int64
}

func (sc *SetCommand) Apply(sm *StateMachine) error {
	sm.state[sc.Key] = int(sc.Value)
	return nil
}
