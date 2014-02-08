package main

import "testing"

var game *Game

func setup() {
	game = NewGame(nil)
}

func assertDirections(t *testing.T, dirs []Direction) {
	for i, _ := range game.snake {
		if game.snake[i].nextDir != dirs[i] {
			t.Errorf("Wrong direction part %d, expected %d but was %d", i, dirs[i], game.snake[i].nextDir)
		}
	}
}

func assertPositions(t *testing.T, x, y int) {
	pos := &game.snake[0].pos
	if pos.X != x || pos.Y != y {
		t.Errorf("Wrong position, expected (%d,%d), but was (%d,%d)", x, y, pos.X, pos.Y)
	}
}

func assertSize(t *testing.T, size int) {
	if len(game.snake) != size {
		t.Errorf("Wrong size, expected %d, but was %s", size, len(game.snake))
	}
}

func assertBlock(t *testing.T, x, y int, bt BlockType) {
	if game.grid[x][y] != bt {
		t.Errorf("Wrong block at (%d,%d), expected %d but was %d", x, y, bt, game.grid[x][y])
	}
}

func TestSnakeInit(t *testing.T) {
	setup()

	assertDirections(t, []Direction{START_DIR, START_DIR, START_DIR, START_DIR})
	assertPositions(t, START_X, START_Y)
	assertSize(t, START_LENGTH)
}

func TestNextDir(t *testing.T) {
	setup()

	nextDir(game, UP)

	assertDirections(t, []Direction{UP, RIGHT, RIGHT, RIGHT})
	assertPositions(t, START_X, START_Y)
}

func TestMoveSnake(t *testing.T) {
	setup()

	moveSnake(game)

	assertDirections(t, []Direction{RIGHT, RIGHT, RIGHT, RIGHT})
	assertPositions(t, START_X+1, START_Y)
}

func TestMoveSnake_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game)

	assertDirections(t, []Direction{UP, UP, RIGHT, RIGHT})
	assertPositions(t, START_X, START_Y-1)
}

func TestMoveSnake_2_afterNextDir(t *testing.T) {
	setup()
	nextDir(game, UP)

	moveSnake(game)
	moveSnake(game)

	assertDirections(t, []Direction{UP, UP, UP, RIGHT})
	assertPositions(t, START_X, START_Y-2)
}

func TestMoveSnake_off_x_limits_pops_next_side(t *testing.T) {
	setup()

	// move snake one block off the right screen limits
	for i := 0; i < NB_BLOCK_WIDTH-START_X; i++ {
		moveSnake(game)
	}

	// assert it has poped to the right side
	assertPositions(t, 0, START_Y)
}

func TestMoveSnake_off_y_limits_pops_next_side(t *testing.T) {
	setup()
	nextDir(game, UP)

	// move snake one block off the right screen limits
	for i := 0; i < NB_BLOCK_HEIGHT-START_Y+1; i++ {
		moveSnake(game)
	}

	// assert it has poped to the right side
	assertPositions(t, START_X, NB_BLOCK_HEIGHT-1)
}

func TestMoveSnake_on_apple(t *testing.T) {
	setup()
	game.grid[START_X+1][START_Y] = APPLE

	moveSnake(game)

	// assert apple disapear
	assertBlock(t, START_X+1, START_Y, EMPTY)
	// assert size increased
	assertSize(t, START_LENGTH+1)
}
