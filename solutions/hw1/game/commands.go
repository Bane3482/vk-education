package main

func (g *Game) LookAround(_ []string) string {
	room := g.Player.Location
	transitions := g.World.GetTransitions(room.ID)

	return room.WithExits(transitions, room.LookDescription(g))
}

func (g *Game) Move(args []string) string {
	transitions := g.World.GetTransitions(g.Player.Location.ID)

	for _, exit := range transitions {
		if exit.Name != args[0] {
			continue
		}
		destination := g.World.Rooms[exit.To]
		if destination == nil {
			break
		}
		if exit.Entry != nil && exit.Entry.Locked {
			return "дверь закрыта"
		}
		g.Player.Location = destination

		transitions = g.World.GetTransitions(destination.ID)

		return destination.WithExits(transitions, destination.EntryDescription)
	}
	return "нет пути в " + args[0]
}

func (g *Game) Wear(args []string) string {
	item, furniture := g.Player.Location.FindItem(args[0])
	if item == nil {
		return "нет такого"
	}
	if !item.Wearable {
		return "нельзя надеть: " + item.Name
	}
	g.Player.Equip(item)
	furniture.RemoveItem(item)
	return "вы надели: " + item.Name
}

func (g *Game) Take(args []string) string {
	item, furniture := g.Player.Location.FindItem(args[0])
	if item == nil {
		return "нет такого"
	}
	if !g.Player.CanStore() {
		return "некуда класть"
	}
	g.Player.AddItem(item)
	furniture.RemoveItem(item)
	return "предмет добавлен в инвентарь: " + item.Name
}

func (g *Game) Use(args []string) string {
	item := g.Player.GetItem(args[0])
	if item == nil {
		return "нет предмета в инвентаре - " + args[0]
	}
	target := g.Player.Location.Targets[args[1]]
	if target != nil {
		if rule := target.UseRules[item.ID]; rule != nil {
			return rule(g)
		}
	}
	return "не к чему применить"
}
