package rooms

type Item interface {
}

type Room struct {
	ID          string
	Description string

	Items map[string]Item
}

func NewRoom(items map[string]Item) *Room {
	return &Room{
		Items: items,
	}
}

func (r *Room) GetItem(id string) Item {
	return r.Items[id]
}
