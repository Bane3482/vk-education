package session

type (
	Player interface {
	}

	World interface {
	}
)

type GameSession struct {
	player Player
	world  World
}

func NewGameSession(player Player, world World) *GameSession {
	return &GameSession{player: player, world: world}
}
