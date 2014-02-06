package main

import "testing"

var game *Game

func setup() {
	game = NewGame(nil)
}

func TestSnakeBoot(t *testing.T) {
	setup()

	assertSnake(t, game.snake, []Direction{RIGHT, RIGHT, RIGHT, RIGHT})
}

func TestNextDir(t *testing.T) {
	setup()

	nextDir(game, UP)

	assertSnake(t, game.snake, []Direction{UP, RIGHT, RIGHT, RIGHT})
}

func TestMoveSnake(t *testing.T) {
	setup()

	moveSnake(game.snake)

	assertSnake(t, game.snake, []Direction{RIGHT, RIGHT, RIGHT, RIGHT})
}

func TestMoveSnake_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game.snake)

	assertSnake(t, game.snake, []Direction{UP, UP, RIGHT, RIGHT})
}

func TestMoveSnake_2_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game.snake)
	moveSnake(game.snake)

	assertSnake(t, game.snake, []Direction{UP, UP, UP, RIGHT})
}

func assertSnake(t *testing.T, snake Snake, dirs []Direction) {
	for i, _ := range snake {
		if snake[i].nextDir != dirs[i] {
			t.Errorf("Wrong direction part %d, expected %d but was %d", i, dirs[i], snake[i].nextDir)
		}
	}
}
