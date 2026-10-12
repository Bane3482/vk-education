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
	Name             string
	Locked           bool
	OpenDescription  string
	CloseDescription string
}

type Interaction struct {
	UseRules map[string]func(*Game) string
}

func CreateEntryInteraction(itemID string, entry *Entry) *Interaction {
	return &Interaction{
		UseRules: map[string]func(*Game) string{
			itemID: func(_ *Game) string {
				if entry.Locked {
					entry.Locked = false
					return entry.OpenDescription
				}

				entry.Locked = true
				return entry.CloseDescription
			},
		},
	}
}

func (w *World) GetTransitions(id string) []*Transition {
	if transitions, ok := w.Transitions[id]; ok {
		return transitions
	}

	return make([]*Transition, 0)
}
