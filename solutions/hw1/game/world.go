package main

type Transition struct {
	Name  string
	To    string
	Entry *Entry
}

type World struct {
	Rooms       map[string]*Room
	Transitions map[string][]*Transition
}

type Entry struct {
	Locked bool
}

type Interaction struct {
	UseRules map[string]func(*Game) string
}

func (w *World) GetTransitions(id string) []*Transition {
	if transitions, ok := w.Transitions[id]; ok {
		return transitions
	}

	return make([]*Transition, 0)
}
