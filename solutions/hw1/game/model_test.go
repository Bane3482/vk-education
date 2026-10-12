package main

import "testing"

// A different world exercises the same commands without changing their code.
func TestModelWorldIsConfigurable(t *testing.T) {
	bag := &Item{ID: "bag", Name: "сумка", Wearable: true, ProvidesStorage: true}
	token := &Item{ID: "token", Name: "жетон"}
	note := &Item{ID: "note", Name: "записка"}
	gate := &Entry{Name: "калитка", Locked: true, OpenDescription: "калитка открыта", CloseDescription: "калитка закрыта"}
	workshop := &Room{
		ID: "workshop", EntryDescription: "мастерская", EmptyDescription: "полки пусты",
		Furniture: []*Furniture{{Notation: "на полке", Items: []*Item{bag, token, note}}},
		Targets: map[string]*Interaction{
			gate.Name: CreateEntryInteraction(token.ID, gate),
		},
	}

	garden := &Room{
		ID: "green-space", EntryDescription: "цветущий сад",
	}

	transitions := map[string][]*Transition{
		workshop.ID: []*Transition{{Name: "сад", To: garden.ID, Entry: gate}},
		garden.ID:   []*Transition{{Name: "назад", To: workshop.ID, Entry: gate}},
	}

	g := NewGame(
		&World{
			Rooms:       map[string]*Room{workshop.ID: workshop, garden.ID: garden},
			Transitions: transitions,
		}, NewPlayer(workshop),
	)

	steps := []additionalGameStep{
		{"осмотреться", "на полке: сумка, жетон, записка. можно пройти - сад"},
		{"идти сад", "дверь закрыта"},
		{"взять жетон", "некуда класть"},
		{"надеть записка", "нельзя надеть: записка"},
		{"осмотреться", "на полке: сумка, жетон, записка. можно пройти - сад"},
		{" \tнадеть   сумка  ", "вы надели: сумка"},
		{"взять жетон", "предмет добавлен в инвентарь: жетон"},
		{"осмотреться", "на полке: записка. можно пройти - сад"},
		{"применить жетон калитка", "калитка открыта"},
		{"идти сад", "цветущий сад. можно пройти - назад"},
		{"осмотреться", "цветущий сад. можно пройти - назад"},
		{"применить жетон калитка", "не к чему применить"},
		{"идти назад", "мастерская. можно пройти - сад"},
		{"взять записка", "предмет добавлен в инвентарь: записка"},
		{"осмотреться", "полки пусты. можно пройти - сад"},
	}
	for i, step := range steps {
		if got := g.HandleCommand(step.command); got != step.want {
			t.Fatalf("step %d, command %q:\n got: %q\nwant: %q", i+1, step.command, got, step.want)
		}
	}

	// Surplus arguments must not accidentally execute an otherwise valid action.
	for _, command := range []string{
		"осмотреться лишнее", "идти сад лишнее", "надеть сумка лишнее",
		"взять жетон лишнее", "применить жетон калитка лишнее",
	} {
		if got := g.HandleCommand(command); got != "неизвестная команда" {
			t.Fatalf("command %q: got %q", command, got)
		}
	}
	if g.Player.Location != workshop || !g.Player.HasItem(token.ID) || gate.Locked {
		t.Fatal("malformed commands changed location, inventory or door state")
	}

	// A room with no exits should not produce an empty list suffix.
	transitions[garden.ID] = nil
	if got := g.HandleCommand("идти сад"); got != "цветущий сад" {
		t.Fatalf("room without exits: got %q", got)
	}
}

func TestModelGamesHaveIndependentState(t *testing.T) {
	newGame := func() *Game {
		room := &Room{ID: "start", EntryDescription: "комната", EmptyDescription: "пусто",
			Furniture: []*Furniture{{Notation: "на полке", Items: []*Item{
				{ID: "bag", Name: "сумка", Wearable: true, ProvidesStorage: true},
				{ID: "token", Name: "жетон"},
			}}},
		}
		return NewGame(&World{Rooms: map[string]*Room{room.ID: room}}, NewPlayer(room))
	}
	first, second := newGame(), newGame()
	for _, step := range []additionalGameStep{
		{"надеть сумка", "вы надели: сумка"},
		{"взять жетон", "предмет добавлен в инвентарь: жетон"},
		{"осмотреться", "пусто"},
	} {
		if got := first.HandleCommand(step.command); got != step.want {
			t.Fatalf("first game, command %q: got %q, want %q", step.command, got, step.want)
		}
	}
	if got := second.HandleCommand("осмотреться"); got != "на полке: сумка, жетон" {
		t.Fatalf("second game's items changed: %q", got)
	}
	if got := second.HandleCommand("взять жетон"); got != "некуда класть" {
		t.Fatalf("second game inherited equipment: %q", got)
	}
}
