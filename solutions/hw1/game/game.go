package main

import "strings"

type Command struct {
	ArgCount int
	Execute  func([]string) string
}

type Game struct {
	World    *World
	Player   *Player
	Commands map[string]Command
}

func NewGame(world *World, player *Player) *Game {
	g := &Game{World: world, Player: player}
	g.Commands = map[string]Command{
		"осмотреться": {ArgCount: 0, Execute: g.LookAround},
		"идти":        {ArgCount: 1, Execute: g.Move},
		"надеть":      {ArgCount: 1, Execute: g.Wear},
		"взять":       {ArgCount: 1, Execute: g.Take},
		"применить":   {ArgCount: 2, Execute: g.Use},
	}
	return g
}

func (g *Game) HandleCommand(input string) string {
	words := strings.Fields(input)
	if len(words) == 0 {
		return "неизвестная команда"
	}
	command, ok := g.Commands[words[0]]
	if !ok || len(words)-1 != command.ArgCount {
		return "неизвестная команда"
	}
	return command.Execute(words[1:])
}
