package main

import (
	"testing"
)

var g *Game

func setup() {
	g = NewGame()
}

func TestCollision(t *testing.T) {
	c1 := Candy{x: 10, y: 20}
	c2 := Candy{x: 5, y: 10}

	collision := collide(c1, c2)

	if !collision {
		t.Errorf("(%d,%d) and (%d,%d) should collide", c1.x, c1.y, c2.x, c2.y)
	}
}

func TestCollisionColumn(t *testing.T) {
	col := Column{candys: []Candy{Candy{x: 0, y: 0}, Candy{x: 0, y: BlockSize}}}

	collision := collideColumnInd(0, col)

	if collision {
		t.Errorf("(0,0) and (0,%d) should not collide", BlockSize)
	}
}

func TestNoCollision(t *testing.T) {
	c1 := Candy{x: 100, y: 20}
	c2 := Candy{x: 5, y: 10}

	collision := collide(c1, c2)

	if collision {
		t.Errorf("(%d,%d) and (%d,%d) should not collide", c1.x, c1.y, c2.x, c2.y)
	}

}

func TestPopulateDropZone(t *testing.T) {
	setup()

	g.populateDropZone()

	for _, col := range g.columns {
		assertNbCandy(t, col, 1)
		assertNotEmpty(t, col.candys[0])
	}
}

func TestMove(t *testing.T) {
	setup()
	g.populateDropZone()
	g.applyVectors()

	moving := g.move()

	assertNbCandy(t, g.columns[0], 1)
	assertY(t, g.columns[0].candys[0], 1)
	if !moving {
		t.Error("Wrong move state, should still moving")
	}
}

func TestGenerateCandy(t *testing.T) {
	setup()

	candy := g.newCandy()

	assertNotEmpty(t, candy)
}

func TestCheckLineNoMatch(t *testing.T) {
	setup()
	line := []CandyType{RedCandy, RedCandy, YellowCandy, YellowCandy, BlueCandy, BlueCandy}

	match := checkLine(line)

	assertMatch(t, match, 0, 0)
}

func TestCheckLineMatch3(t *testing.T) {
	setup()
	line := []CandyType{RedCandy, BlueCandy, BlueCandy, BlueCandy, YellowCandy, GreenCandy, GreenCandy}

	match := checkLine(line)

	assertMatch(t, match, 1, Match3)
}

func TestCheckLineMatch4(t *testing.T) {
	setup()
	line := []CandyType{RedCandy, GreenCandy, BlueCandy, BlueCandy, BlueCandy, BlueCandy, GreenCandy, GreenCandy}

	match := checkLine(line)

	assertMatch(t, match, 2, Match4)
}

func TestCheckLineMatch5(t *testing.T) {
	setup()
	line := []CandyType{RedCandy, BlueCandy, BlueCandy, BlueCandy, BlueCandy, BlueCandy, GreenCandy}

	match := checkLine(line)

	assertMatch(t, match, 1, Match5)
}

func TestCheckGridNoMatch(t *testing.T) {
	setup()
	candys := make([][]CandyType, 4)
	candys[0] = []CandyType{RedCandy, GreenCandy, YellowCandy, BlueCandy}
	candys[1] = []CandyType{RedCandy, GreenCandy, YellowCandy, BlueCandy}
	candys[2] = []CandyType{RedCandy, GreenCandy, YellowCandy, BlueCandy}
	candys[3] = []CandyType{RedCandy, GreenCandy, YellowCandy, BlueCandy}

	matches := checkGrid(candys)

	for _, match := range matches {
		assertMatch(t, match, 0, 0)
	}

}

func TestApplyVector(t *testing.T) {
	setup()
	g.populateDropZone()

	g.applyVectors()

	for _, col := range g.columns {
		assertVector(t, col.candys[0].v, 1)
	}
}

func assertMatch(t *testing.T, match Match, start, length int) {
	if match.start != start && match.length != length {
		t.Errorf("Wrong match, expected start=%d, length=%d, but was (%d,%d)", start, length, match.start, match.length)
	}
}

func assertVector(t *testing.T, v, expected int) {
	if v != expected {
		t.Errorf("Wrong candy vector, expected %d but was %d", expected, v)
	}
}

func assertCandy(t *testing.T, c Candy, expected CandyType) {
	if c._type != expected {
		t.Errorf("Wrong candy type, expected %d but was %d", expected, c._type)
	}
}

func assertNotEmpty(t *testing.T, c Candy) {
	if c._type == EmptyCandy {
		t.Errorf("Wrong candy type, expected not empty")
	}
}

func assertY(t *testing.T, c Candy, y int) {
	if c.y != y {
		t.Errorf("Wrong y, expected %d but was %d", y, c.y)
	}
}

func assertNbCandy(t *testing.T, col Column, nb int) {
	if len(col.candys) != nb {
		t.Fatalf("Wrong number of candy in column, expected %d but was %d", nb, len(col.candys))
	}
}
