package main

import "testing"

var game *Game

func setup() {
	game = NewGame(nil)
}

func assertSnake(t *testing.T, snake Snake, dirs []Direction) {
	for i, _ := range snake {
		if snake[i].nextDir != dirs[i] {
			t.Errorf("Wrong direction part %d, expected %d but was %d", i, dirs[i], snake[i].nextDir)
		}
	}
}

func assertSnakePosition(t *testing.T, snake Snake, x, y int) {
	pos := &snake[0].pos
	if pos.X != x || pos.Y != y {
		t.Errorf("Wrong position, expected (%d,%d), but was (%d,%d)", x, y, pos.X, pos.Y)
	}
}

func TestSnakeBoot(t *testing.T) {
	setup()

	assertSnake(t, game.snake, []Direction{RIGHT, RIGHT, RIGHT, RIGHT})
	assertSnakePosition(t, game.snake, START_X, START_Y)
}

func TestNextDir(t *testing.T) {
	setup()

	nextDir(game, UP)

	assertSnake(t, game.snake, []Direction{UP, RIGHT, RIGHT, RIGHT})
	assertSnakePosition(t, game.snake, START_X, START_Y)
}

func TestMoveSnake(t *testing.T) {
	setup()

	moveSnake(game)

	assertSnake(t, game.snake, []Direction{RIGHT, RIGHT, RIGHT, RIGHT})
	assertSnakePosition(t, game.snake, START_X+1, START_Y)
}

func TestMoveSnake_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game)

	assertSnake(t, game.snake, []Direction{UP, UP, RIGHT, RIGHT})
	assertSnakePosition(t, game.snake, START_X, START_Y-1)
}

func TestMoveSnake_2_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game)
	moveSnake(game)

	assertSnake(t, game.snake, []Direction{UP, UP, UP, RIGHT})
	assertSnakePosition(t, game.snake, START_X, START_Y-2)
}

func TestMoveSnake_off_screen_limits_pops_next_side(t *testing.T) {
	setup()

	// move snake one block off the right screen limits
	for i := 0; i < NB_BLOCK_WIDTH/2; i++ {
		moveSnake(game)
	}

	// assert it has poped to the right side
	assertSnakePosition(t, game.snake, 0, START_Y)
}
