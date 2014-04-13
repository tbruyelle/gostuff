package main

import "testing"

var game *Game

func setup() {
	game = NewGame(nil)
}

func assertSnakeRate(t *testing.T, snakeRate float32) {
	if game.snakeRate != snakeRate {
		t.Errorf("Wrong snake rate, expected %f but was %f", snakeRate, game.snakeRate)
	}
}

func assertHeadPosition(t *testing.T, x, y int) {
	pos := &game.Snake[0].Pos
	if pos.X != x || pos.Y != y {
		t.Errorf("Wrong position, expected (%d,%d), but was (%d,%d)", x, y, pos.X, pos.Y)
	}
}

func assertSnakeDirections(t *testing.T, dirs ...Direction) {
	pos := game.Snake[0].Pos
	for i, dir := range dirs {
		movePos(dir, &pos)
		cpos := &game.Snake[i+1].Pos
		if cpos.X != pos.X || cpos.Y != pos.Y {
			t.Errorf("Wrong queue position, expected (%d,%d), but was (%d,%d)", pos.X, pos.Y, cpos.X, cpos.Y)
		}
	}
}

func assertSize(t *testing.T, size int) {
	if len(game.Snake) != size {
		t.Errorf("Wrong size, expected %d, but was %s", size, len(game.Snake))
	}
}

func assertBlock(t *testing.T, x, y int, bt BlockType) {
	if game.Grid[x][y] != bt {
		t.Errorf("Wrong block at (%d,%d), expected %d but was %d", x, y, bt, game.Grid[x][y])
	}
}

func assertRunning(t *testing.T, expected bool) {
	var running bool
	select {
	case <-game.EndLoop:
		running = false
	default:
		running = true
		// OK
	}
	if expected != running {
		t.Errorf("Wrong loop status, expected %t but was %t", expected, running)
	}
}

func TestSnakeInit(t *testing.T) {
	setup()

	assertRunning(t, true)
	assertHeadPosition(t, START_X, START_Y)
	assertSnakeDirections(t, -START_DIR, -START_DIR, -START_DIR)
	assertSize(t, START_LENGTH)
}

func TestCommand(t *testing.T) {
	setup()

	game.Command(UP)

	assertRunning(t, true)
	assertSnakeDirections(t, -START_DIR, -START_DIR, -START_DIR)
	assertHeadPosition(t, START_X, START_Y)
}

func TestMoveSnake(t *testing.T) {
	setup()

	game.Tick()

	assertRunning(t, true)
	assertHeadPosition(t, START_X+1, START_Y)
	assertSnakeDirections(t, -START_DIR, -START_DIR, -START_DIR)
}

func TestMoveSnake_afterCommand(t *testing.T) {
	setup()
	game.Command(UP)

	game.Tick()

	assertRunning(t, true)
	assertHeadPosition(t, START_X, START_Y-1)
	assertSnakeDirections(t, DOWN, -START_DIR, -START_DIR)
}

func TestMoveSnake_2_afterCommand(t *testing.T) {
	setup()
	game.Command(UP)

	game.Tick()
	game.Tick()

	assertRunning(t, true)
	assertHeadPosition(t, START_X, START_Y-2)
	assertSnakeDirections(t, DOWN, DOWN, -START_DIR)
}

func TestMoveSnake_off_x_limits_pops_next_side(t *testing.T) {
	setup()

	// move snake one block off the right screen limits
	for i := 0; i < NB_BLOCK_WIDTH-START_X; i++ {
		game.Tick()
	}

	assertRunning(t, true)
	// assert it has poped to the right side
	assertHeadPosition(t, 0, START_Y)
}

func TestMoveSnake_off_y_limits_pops_next_side(t *testing.T) {
	setup()
	game.Command(UP)

	// move snake one block off the right screen limits
	for i := 0; i < NB_BLOCK_HEIGHT-START_Y+1; i++ {
		game.Tick()
	}

	assertRunning(t, true)
	// assert it has poped to the right side
	assertHeadPosition(t, START_X, NB_BLOCK_HEIGHT-1)
}

func TestMoveSnake_on_apple(t *testing.T) {
	setup()
	game.Grid[START_X+1][START_Y] = APPLE

	game.Tick()

	assertRunning(t, true)
	// assert apple disapear
	assertBlock(t, START_X+1, START_Y, EMPTY)
	// assert size increased
	assertSize(t, START_LENGTH+1)
}

func TestSnakeCollision(t *testing.T) {
	setup()
	// increase size because the snake must have at least a size of 5
	// to make collision possible
	game.grow()

	game.Command(UP)
	game.Tick()
	game.Command(-START_DIR)
	game.Tick()
	game.Command(DOWN)
	game.Tick()

	assertRunning(t, false)
}

func TestCommandBack_is_ignored(t *testing.T) {
	setup()

	game.Command(UP) // prevent a bug when nextDir is updated before the tick
	game.Command(-START_DIR)
	game.Tick()

	assertRunning(t, true)
	assertHeadPosition(t, START_X, START_Y-1)
}

func TestEatApple_increase_snake_rate(t *testing.T) {
	setup()
	game.Grid[START_X+1][START_Y] = APPLE

	game.Tick()

	assertRunning(t, true)
	assertSnakeRate(t, START_SNAKE_RATE+1)
}
