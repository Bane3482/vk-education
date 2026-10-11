package items

type Backpack struct {
	Item

	inventory map[string]*Item
}
