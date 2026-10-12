package main

import "testing"

const (
	additionalKitchenInitial = "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"
	additionalKitchenReady   = "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
	additionalKitchenEntry   = "кухня, ничего интересного. можно пройти - коридор"
	additionalCorridor       = "ничего интересного. можно пройти - кухня, комната, улица"
	additionalRoomEntry      = "ты в своей комнате. можно пройти - коридор"
	additionalRoomFull       = "на столе: ключи, конспекты, на стуле: рюкзак. можно пройти - коридор"
	additionalRoomTable      = "на столе: ключи, конспекты. можно пройти - коридор"
	additionalRoomEmpty      = "пустая комната. можно пройти - коридор"
	additionalStreet         = "на улице весна. можно пройти - домой"
)

type additionalGameStep struct {
	command string
	want    string
}

func additionalRunSteps(t *testing.T, steps ...additionalGameStep) {
	t.Helper()
	for i, step := range steps {
		if got := handleCommand(step.command); got != step.want {
			t.Fatalf("step %d, command %q:\n got: %q\nwant: %q", i+1, step.command, got, step.want)
		}
	}
}

func additionalEnterRoom(t *testing.T) {
	t.Helper()
	additionalRunSteps(t,
		additionalGameStep{"идти коридор", additionalCorridor},
		additionalGameStep{"идти комната", additionalRoomEntry},
	)
}

func additionalCollectItems(t *testing.T) {
	t.Helper()
	additionalEnterRoom(t)
	additionalRunSteps(t,
		additionalGameStep{"надеть рюкзак", "вы надели: рюкзак"},
		additionalGameStep{"взять ключи", "предмет добавлен в инвентарь: ключи"},
		additionalGameStep{"взять конспекты", "предмет добавлен в инвентарь: конспекты"},
	)
}

// These tests use only initGame and handleCommand, so they do not depend on
// the internal architecture. Do not use t.Parallel: the game has global state.
func TestAdditionalGameScenarios(t *testing.T) {
	cases := []struct {
		name  string
		steps []additionalGameStep
	}{
		{
			name: "failed_movement_keeps_location",
			steps: []additionalGameStep{
				{"идти комната", "нет пути в комната"},
				{"осмотреться", additionalKitchenInitial},
				{"идти несуществующая", "нет пути в несуществующая"},
				{"осмотреться", additionalKitchenInitial},
				{"идти коридор", additionalCorridor},
				{"идти улица", "дверь закрыта"},
				{"осмотреться", additionalCorridor},
				{"идти комната", additionalRoomEntry},
				{"идти кухня", "нет пути в кухня"},
				{"осмотреться", additionalRoomFull},
			},
		},
		{
			name: "failed_pickup_keeps_items",
			steps: []additionalGameStep{
				{"идти коридор", additionalCorridor},
				{"идти комната", additionalRoomEntry},
				{"взять ключи", "некуда класть"},
				{"взять конспекты", "некуда класть"},
				{"осмотреться", additionalRoomFull},
				{"надеть рюкзак", "вы надели: рюкзак"},
				{"осмотреться", additionalRoomTable},
				{"взять телефон", "нет такого"},
				{"осмотреться", additionalRoomTable},
				{"взять ключи", "предмет добавлен в инвентарь: ключи"},
				{"взять конспекты", "предмет добавлен в инвентарь: конспекты"},
				{"осмотреться", additionalRoomEmpty},
			},
		},
		{
			name: "items_can_be_collected_in_reverse_order",
			steps: []additionalGameStep{
				{"идти коридор", additionalCorridor},
				{"идти комната", additionalRoomEntry},
				{"надеть рюкзак", "вы надели: рюкзак"},
				{"взять конспекты", "предмет добавлен в инвентарь: конспекты"},
				{"осмотреться", "на столе: ключи. можно пройти - коридор"},
				{"взять конспекты", "нет такого"},
				{"взять ключи", "предмет добавлен в инвентарь: ключи"},
				{"взять ключи", "нет такого"},
				{"осмотреться", additionalRoomEmpty},
				{"идти коридор", additionalCorridor},
				{"идти кухня", additionalKitchenEntry},
				{"осмотреться", additionalKitchenReady},
			},
		},
		{
			name: "backpack_and_inventory_survive_movement",
			steps: []additionalGameStep{
				{"идти коридор", additionalCorridor},
				{"идти комната", additionalRoomEntry},
				{"надеть рюкзак", "вы надели: рюкзак"},
				{"взять ключи", "предмет добавлен в инвентарь: ключи"},
				{"идти коридор", additionalCorridor},
				{"идти кухня", additionalKitchenEntry},
				{"взять конспекты", "нет такого"},
				{"идти коридор", additionalCorridor},
				{"идти комната", additionalRoomEntry},
				{"осмотреться", "на столе: конспекты. можно пройти - коридор"},
				{"взять конспекты", "предмет добавлен в инвентарь: конспекты"},
				{"осмотреться", additionalRoomEmpty},
				{"идти коридор", additionalCorridor},
				{"применить ключи дверь", "дверь открыта"},
				{"идти улица", additionalStreet},
			},
		},
		{
			name: "unknown_commands_keep_state",
			steps: []additionalGameStep{
				{"завтракать", "неизвестная команда"},
				{"осмотреться", additionalKitchenInitial},
				{"идтикомната", "неизвестная команда"},
				{"осмотреться", additionalKitchenInitial},
				{"идти коридор", additionalCorridor},
				{"идти комната", additionalRoomEntry},
				{"надеть рюкзак", "вы надели: рюкзак"},
				{"взятьключи", "неизвестная команда"},
				{"осмотреться", additionalRoomTable},
				{"взять ключи", "предмет добавлен в инвентарь: ключи"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initGame()
			additionalRunSteps(t, tc.steps...)
		})
	}
}

func TestAdditionalDoorAndUsageState(t *testing.T) {
	initGame()
	additionalCollectItems(t)
	additionalRunSteps(t,
		// A target in another room must not be accessible.
		additionalGameStep{"применить ключи дверь", "не к чему применить"},
		additionalGameStep{"осмотреться", additionalRoomEmpty},
		additionalGameStep{"идти коридор", additionalCorridor},
		// Having keys is not enough: the player must explicitly use them.
		additionalGameStep{"идти улица", "дверь закрыта"},
		additionalGameStep{"применить телефон дверь", "нет предмета в инвентаре - телефон"},
		additionalGameStep{"применить ключи шкаф", "не к чему применить"},
		additionalGameStep{"применить конспекты дверь", "не к чему применить"},
		additionalGameStep{"идти улица", "дверь закрыта"},
		additionalGameStep{"применить ключи дверь", "дверь открыта"},
		// Using the keys must not consume them.
		additionalGameStep{"применить ключи шкаф", "не к чему применить"},
		// Revisiting a room must not restore its items or relock the door.
		additionalGameStep{"идти комната", additionalRoomEntry},
		additionalGameStep{"осмотреться", additionalRoomEmpty},
		additionalGameStep{"идти коридор", additionalCorridor},
		additionalGameStep{"идти улица", additionalStreet},
	)
}

func TestAdditionalInitGameResetsAllState(t *testing.T) {
	initGame()
	additionalCollectItems(t)
	additionalRunSteps(t,
		additionalGameStep{"идти коридор", additionalCorridor},
		additionalGameStep{"применить ключи дверь", "дверь открыта"},
		additionalGameStep{"идти улица", additionalStreet},
	)

	// Reset both a completed game and a partially modified game.
	for reset := 1; reset <= 2; reset++ {
		initGame()
		additionalRunSteps(t,
			additionalGameStep{"осмотреться", additionalKitchenInitial},
			additionalGameStep{"идти коридор", additionalCorridor},
			additionalGameStep{"идти улица", "дверь закрыта"},
			additionalGameStep{"применить ключи дверь", "нет предмета в инвентаре - ключи"},
			additionalGameStep{"применить конспекты дверь", "нет предмета в инвентаре - конспекты"},
			additionalGameStep{"идти комната", additionalRoomEntry},
			additionalGameStep{"осмотреться", additionalRoomFull},
			additionalGameStep{"взять ключи", "некуда класть"},
			additionalGameStep{"надеть рюкзак", "вы надели: рюкзак"},
			additionalGameStep{"взять ключи", "предмет добавлен в инвентарь: ключи"},
		)
	}
}

func TestAdditionalLookAroundDoesNotChangeState(t *testing.T) {
	initGame()
	for i := 0; i < 20; i++ {
		additionalRunSteps(t, additionalGameStep{"осмотреться", additionalKitchenInitial})
	}
	additionalEnterRoom(t)
	for i := 0; i < 20; i++ {
		additionalRunSteps(t, additionalGameStep{"осмотреться", additionalRoomFull})
	}
	additionalRunSteps(t,
		additionalGameStep{"надеть рюкзак", "вы надели: рюкзак"},
		additionalGameStep{"взять ключи", "предмет добавлен в инвентарь: ключи"},
		additionalGameStep{"взять конспекты", "предмет добавлен в инвентарь: конспекты"},
	)
	for i := 0; i < 20; i++ {
		additionalRunSteps(t, additionalGameStep{"осмотреться", additionalRoomEmpty})
	}
}

// Robustness beyond the original scenarios: malformed input must not panic
// or mutate the game. The exact error text is deliberately not prescribed.
func TestAdditionalMalformedCommandsKeepState(t *testing.T) {
	commands := []string{"", "   ", "\t\n", "идти", "взять", "надеть", "применить", "применить ключи"}
	for _, command := range commands {
		t.Run("command="+command, func(t *testing.T) {
			initGame()
			additionalCollectItems(t)
			additionalRunSteps(t, additionalGameStep{"идти коридор", additionalCorridor})
			_ = handleCommand(command)
			additionalRunSteps(t,
				additionalGameStep{"осмотреться", additionalCorridor},
				additionalGameStep{"идти улица", "дверь закрыта"},
				additionalGameStep{"применить ключи дверь", "дверь открыта"},
				additionalGameStep{"идти комната", additionalRoomEntry},
				additionalGameStep{"осмотреться", additionalRoomEmpty},
				additionalGameStep{"идти коридор", additionalCorridor},
				additionalGameStep{"идти кухня", additionalKitchenEntry},
				additionalGameStep{"осмотреться", additionalKitchenReady},
			)
		})
	}
}
