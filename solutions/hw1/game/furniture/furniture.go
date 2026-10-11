package furniture

type Item interface {
}

type Furniture struct {
	ID          string
	Description string

	Items []Item
}

func NewFurniture(id string, description string, items []Item) *Furniture {
	return &Furniture{
		ID:          id,
		Description: description,
		Items:       items,
	}
}
