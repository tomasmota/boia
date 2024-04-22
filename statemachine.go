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

//
// func () Apply(entry LogEntry) error {
// 	sm.state[entry.Command.Key] = entry.Command.Value
// 	return nil
// }
//
// func (sm *StateMachine) Get(key string) (int, error) {
// 	v, ok := sm.state[key]
// 	if !ok {
// 		return 0, fmt.Errorf(`no entry with key (%s) found in state machine`, key)
// 	}
// 	return v, nil
// }
