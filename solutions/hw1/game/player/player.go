package player

type (
	Item interface {
	}

	Room interface {
	}

	Backpack interface {
		GetItem(id string) Item
	}
)

type player struct {
	nickname string

	location  Room
	inventory Backpack
}

func NewPlayer(nickname string, location Room, inventory Backpack) *player {
	return &player{nickname: nickname, location: location, inventory: inventory}
}
