package world

type (
	Room interface {
	}

	Transition interface {
	}
)

type World struct {
	rooms       []Room
	transitions []Transition
}

func NewWorld(rooms []Room, transitions []Transition) *World {
	return &World{rooms: rooms, transitions: transitions}
}
