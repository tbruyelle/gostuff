package main

import "testing"

var game *Game

func setup() {
	game = NewGame(nil)
}

func assertDirections(t *testing.T, g *Game, dirs []Direction) {
	for i, _ := range g.snake {
		if g.snake[i].nextDir != dirs[i] {
			t.Errorf("Wrong direction part %d, expected %d but was %d", i, dirs[i], g.snake[i].nextDir)
		}
	}
}

func assertPositions(t *testing.T, g *Game, x, y int) {
	pos := &g.snake[0].pos
	if pos.X != x || pos.Y != y {
		t.Errorf("Wrong position, expected (%d,%d), but was (%d,%d)", x, y, pos.X, pos.Y)
	}
}

func TestSnakeBoot(t *testing.T) {
	setup()

	assertDirections(t, game, []Direction{RIGHT, RIGHT, RIGHT, RIGHT})
	assertPositions(t, game, START_X, START_Y)
}

func TestNextDir(t *testing.T) {
	setup()

	nextDir(game, UP)

	assertDirections(t, game, []Direction{UP, RIGHT, RIGHT, RIGHT})
	assertPositions(t, game, START_X, START_Y)
}

func TestMoveSnake(t *testing.T) {
	setup()

	moveSnake(game)

	assertDirections(t, game, []Direction{RIGHT, RIGHT, RIGHT, RIGHT})
	assertPositions(t, game, START_X+1, START_Y)
}

func TestMoveSnake_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game)

	assertDirections(t, game, []Direction{UP, UP, RIGHT, RIGHT})
	assertPositions(t, game, START_X, START_Y-1)
}

func TestMoveSnake_2_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game)
	moveSnake(game)

	assertDirections(t, game, []Direction{UP, UP, UP, RIGHT})
	assertPositions(t, game, START_X, START_Y-2)
}

func TestMoveSnake_off_x_limits_pops_next_side(t *testing.T) {
	setup()

	// move snake one block off the right screen limits
	for i := 0; i < NB_BLOCK_WIDTH-START_X; i++ {
		moveSnake(game)
	}

	// assert it has poped to the right side
	assertPositions(t, game, 0, START_Y)
}

func TestMoveSnake_off_y_limits_pops_next_side(t *testing.T) {
	setup()
	nextDir(game, UP)

	// move snake one block off the right screen limits
	for i := 0; i < NB_BLOCK_HEIGHT-START_Y+1; i++ {
		moveSnake(game)
	}

	// assert it has poped to the right side
	assertPositions(t, game, START_X, NB_BLOCK_HEIGHT-1)
}
