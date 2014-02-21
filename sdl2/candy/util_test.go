package main

import "testing"

var g *Game

func setup() {
	g = NewGame()
}

func fillGame() {
	g.populateDropZone()
	g.applyGravity()
	for g.fall() {
		g.populateDropZone()
		g.applyGravity()
	}
	g.populateDropZone()
}

func assertGravity(t *testing.T, c *Candy, expected int) {
	if c.g != expected {
		t.Errorf("Wrong candy vector, expected %d but was %d", expected, c.g)
	}
}

func assertNotEmpty(t *testing.T, c *Candy) {
	if c._type == EmptyCandy {
		t.Errorf("Wrong candy type, expected not empty")
	}
}

func assertXY(t *testing.T, c *Candy, x, y int) {
	if c.y != y || c.x != x {
		t.Errorf("Wrong x,y, expected %d,%d but was %d,%d", x, y, c.x, c.y)
	}
}

func assertY(t *testing.T, c *Candy, y int) {
	if c.y != y {
		t.Errorf("Wrong y, expected %d but was %d", y, c.y)
	}
}

func assertNbCandy(t *testing.T, nb int) {
	if len(g.candys) != nb {
		t.Fatalf("Wrong number of candys, expected %d but was %d", nb, len(g.candys))
	}
}

func assertNear(t *testing.T, near, expected bool, c1, c2 *Candy) {
	if near != expected {
		t.Errorf("Wrong near, expected %t but was %t for candys (%d,%d) and (%d,%d)", near, expected, c1.x, c1.y, c2.x, c2.y)
	}
}

func assertVx(t *testing.T, c *Candy, vx int) {
	if c.vx != vx {
		t.Errorf("Wrong vx, expected %d but was %d", vx, c.vx)
	}
}

func assertCrushes(t *testing.T, cs []*Candy, crush bool, mutation CandyType) {
	for _, c := range cs {
		assertCrush(t, c, crush, mutation)
	}
}

func assertCrush(t *testing.T, c *Candy, crush bool, mutation CandyType) {
	if c.crush != crush {
		t.Fatalf("Wrong matching for %v, expected %t but was %t", c, crush, c.crush)
	}
	if c.mutation != mutation {
		t.Fatalf("Wrong mutation for %v, expected %d but was %d", c, mutation, c.mutation)
	}
}

type C struct {
	d Direction
	t CandyType
}

func generateCandys(cs ...C) []*Candy {
	region := []*Candy{}
	curx, cury := XMin, YMin
	for _, c := range cs {
		switch c.d {
		case Left:
			curx -= BlockSize
		case Right:
			curx += BlockSize
		case Top:
			cury -= BlockSize
		case Bottom:
			cury += BlockSize
		}
		region = append(region, &Candy{x: curx, y: cury, _type: c.t})
	}
	return region
}

func assertCandyType(t *testing.T, ct CandyType, expected CandyType) {
	if ct != expected {
		t.Errorf("Wrong candy type, expected %d but was %d", expected, ct)
	}
}

func assertMatch(t *testing.T, match, expected bool) {
	if match != expected {
		t.Errorf("Wrong match, expected %t but was %t", expected, match)
	}
}

func assertStriped(t *testing.T, c *Candy, expected bool) {
	if c.isStriped() != expected {
		t.Errorf("Wrong type %v, expected striped=%t but was %t", c, expected, c.isStriped())
	}
}
func assertStripedH(t *testing.T, c *Candy, expected bool) {
	if c.isStripedH() != expected {
		t.Errorf("Wrong type %v, expected striped h=%t but was %t", c, expected, c.isStripedH())
	}
}
func assertStripedV(t *testing.T, c *Candy, expected bool) {
	if c.isStripedV() != expected {
		t.Errorf("Wrong type %v, expected striped v=%t but was %t", c, expected, c.isStripedV())
	}
}
