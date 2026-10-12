package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var currentGame *Game

func main() {
	initGame()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(handleCommand(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func initGame() {
	tea := &Item{ID: "tea", Name: "чай"}
	keys := &Item{ID: "keys", Name: "ключи"}
	notes := &Item{ID: "notes", Name: "конспекты"}
	backpack := &Item{ID: "backpack", Name: "рюкзак", Wearable: true, ProvidesStorage: true}

	kitchen := &Room{
		ID: "kitchen", EntryDescription: "кухня, ничего интересного",
		Furniture: []*Furniture{{Notation: "на столе", Items: []*Item{tea}}},
	}
	bedroom := &Room{
		ID: "bedroom", EntryDescription: "ты в своей комнате", EmptyDescription: "пустая комната",
		Furniture: []*Furniture{
			{Notation: "на столе", Items: []*Item{keys, notes}},
			{Notation: "на стуле", Items: []*Item{backpack}},
		},
	}
	corridor := &Room{ID: "corridor", EntryDescription: "ничего интересного"}
	street := &Room{ID: "street", EntryDescription: "на улице весна"}
	door := &Entry{
		Name:             "дверь",
		Locked:           true,
		OpenDescription:  "дверь открыта",
		CloseDescription: "дверь закрыта",
	}

	transitions := map[string][]*Transition{
		kitchen.ID: {
			{Name: "коридор", To: corridor.ID},
		},
		bedroom.ID: {
			{Name: "коридор", To: corridor.ID},
		},
		corridor.ID: {
			{Name: "кухня", To: kitchen.ID},
			{Name: "комната", To: bedroom.ID},
			{Name: "улица", To: street.ID, Entry: door},
		},
		street.ID: {
			{Name: "домой", To: corridor.ID, Entry: door},
		},
	}

	doorTarget := CreateEntryInteraction(keys.ID, door)

	corridor.Targets = map[string]*Interaction{door.Name: doorTarget}
	street.Targets = map[string]*Interaction{door.Name: doorTarget}

	kitchen.Describe = func(g *Game) string {
		parts := []string{"ты находишься на кухне"}
		if contents := kitchen.Contents(); contents != "" {
			parts = append(parts, contents)
		}
		task := "надо собрать рюкзак и идти в универ"
		if g.Player.Equipment[backpack.ID] != nil && g.Player.HasItem(keys.ID) && g.Player.HasItem(notes.ID) {
			task = "надо идти в универ"
		}
		return strings.Join(append(parts, task), ", ")
	}

	world := &World{
		Rooms: map[string]*Room{
			kitchen.ID: kitchen, bedroom.ID: bedroom, corridor.ID: corridor, street.ID: street,
		},
		Transitions: transitions,
	}

	currentGame = NewGame(world, NewPlayer(kitchen))
}

func handleCommand(command string) string {
	return currentGame.HandleCommand(command)
}
