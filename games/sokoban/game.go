package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	currentGus *sdl.Texture
)

var block = sdl.Rect{W: BLOCK_SIZE, H: BLOCK_SIZE}

func play(renderer *sdl.Renderer) {

	game := NewGame(renderer)
	defer game.Destroy()
	loadLevel(1, game)

	// find player
	loopGrid(func(i, j int) {
		if game.grid[i][j] == GUS {
			game.playerPos.X = int32(i)
			game.playerPos.Y = int32(j)
			game.grid[i][j] = EMPTY
		}
	})

	// main loop
	running := true
	for running {
		if game.goals == 0 {
			//win
			running = false
		}
		var event sdl.Event
		sdl.WaitEvent(&event)
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyDownEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				running = false
			case sdl.K_UP:
				movePlayer(UP, game)
			case sdl.K_DOWN:
				movePlayer(DOWN, game)
			case sdl.K_LEFT:
				movePlayer(LEFT, game)
			case sdl.K_RIGHT:
				movePlayer(RIGHT, game)
			}
		}
		renderGame(renderer, game)
	}
}

func movePlayer(dir Direction, g *Game) {
	g.currentPlayer = g.players[dir]
	pos := sdl.Rect{X: g.playerPos.X, Y: g.playerPos.Y}
	// move and check if correct
	if !movePos(dir, &pos, g) {
		return
	}
	if g.grid[pos.X][pos.Y] == BOX {
		// player want to move a box, try its possible
		boxPos := sdl.Rect{X: pos.X, Y: pos.Y}
		if !movePos(dir, &boxPos, g) || g.grid[boxPos.X][boxPos.Y] == BOX {
			return
		}
		// move the box
		g.grid[pos.X][pos.Y] = EMPTY
		if g.grid[boxPos.X][boxPos.Y] == GOAL {
			g.grid[boxPos.X][boxPos.Y] = BOX_OK
			g.goals--
		} else {
			g.grid[boxPos.X][boxPos.Y] = BOX
		}
	}
	// move player ok OK
	g.playerPos.X = pos.X
	g.playerPos.Y = pos.Y
}

func movePos(dir Direction, pos *sdl.Rect, g *Game) bool {
	switch dir {
	case UP:
		pos.Y--
	case DOWN:
		pos.Y++
	case LEFT:
		pos.X--
	case RIGHT:
		pos.X++
	}
	return pos.X >= 0 && pos.X < NB_BLOCK_HEIGHT && pos.Y >= 0 && pos.Y < NB_BLOCK_WIDTH && g.grid[pos.X][pos.Y] != WALL && g.grid[pos.X][pos.Y] != BOX_OK
}

func renderGame(renderer *sdl.Renderer, g *Game) {
	renderThings(renderer, g)
	// show player
	block.X = g.playerPos.X * BLOCK_SIZE
	block.Y = g.playerPos.Y * BLOCK_SIZE
	renderer.Copy(g.currentPlayer, nil, &block)
	renderer.Present()
}

func loadLevel(level int, g *Game) {
	data, err := ioutil.ReadFile(fmt.Sprintf("levels/%d.lvl", level))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load level %d", level)
		os.Exit(1)
	}
	fmt.Printf("Level loaded with %d tiles\n", len(data))
	g.goals = 0
	loopGrid(func(i, j int) {
		block, err := strconv.ParseInt(string(data[i*NB_BLOCK_WIDTH+j]), 10, 10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load level %d, unable to parse", level)
			os.Exit(1)
		}
		g.grid[i][j] = BlockType(block)
		if g.grid[i][j] == GOAL {
			g.goals++
		}
		//fmt.Printf("data=%s map=%x\n", string(data[i*NB_BLOCK_WIDTH+j]), grid[i][j])
	})
}
