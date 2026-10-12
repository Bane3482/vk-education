package main

type Player struct {
	Location  *Room
	Equipment map[string]*Item
	Inventory map[string]*Item
}

func NewPlayer(location *Room) *Player {
	return &Player{
		Location: location, Equipment: make(map[string]*Item), Inventory: make(map[string]*Item),
	}
}

func (p *Player) CanStore() bool {
	for _, item := range p.Equipment {
		if item.ProvidesStorage {
			return true
		}
	}
	return false
}

func (p *Player) HasItem(id string) bool {
	return p.Inventory[id] != nil
}

func (p *Player) GetItem(name string) *Item {
	for _, item := range p.Inventory {
		if item.Name == name {
			return item
		}
	}
	return nil
}

func (p *Player) AddItem(item *Item) {
	p.Inventory[item.ID] = item
}

func (p *Player) Equip(item *Item) {
	p.Equipment[item.ID] = item
}
