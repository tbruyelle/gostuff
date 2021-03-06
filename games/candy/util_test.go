package main

import "runtime/debug"
import "testing"

var g *Game

func setup() {
	g = NewGame()
}

type candyTypeGenMock struct {
	_type CandyType
}

func (c candyTypeGenMock) NewCandyType() CandyType {
	return c._type
}

func fillGame() {
	g.populateDropZone()
	for !g.Update() {
		g.populateDropZone()
	}
	g.populateDropZone()
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
		t.Errorf("Wrong number of candys, expected %d but was %d", nb, len(g.candys))
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

func assertCrushes(t *testing.T, cs []*Candy, crush bool, _type CandyType) {
	for _, c := range cs {
		assertCrush(t, c, crush, _type)
	}
}

// if _type==EmptyCandy the type assertion is ignored
func assertCrush(t *testing.T, c *Candy, crush bool, _type CandyType) {
	_ = debug.Stack()

	if c.crush != crush {
		debug.PrintStack()
		t.Fatalf("Wrong matching for %v, expected %t but was %t", c, crush, c.crush)
	}
	if _type != EmptyCandy && c._type != _type {
		debug.PrintStack()
		t.Fatalf("Wrong _type for %v, expected %d but was %d", c, _type, c._type)
	}
}

func popCandys(tss [][]CandyType) []*Candy {
	var candys []*Candy
	curx, cury := XMin, YMin
	for _, ts := range tss {
		for _, t := range ts {
			candys = append(candys, &Candy{_type: t, x: curx, y: cury, state: NewIdleState()})
			curx += BlockSize
		}
		cury += BlockSize
		curx = XMin
	}
	return candys
}

func assertCandyTypes(t *testing.T, ctss [][]CandyType) {
	curx, cury := XMin, YMin
	for _, cts := range ctss {
		for _, ct := range cts {
			c, ok := findCandy(g.candys, curx, cury)
			if ct == EmptyCandy {
				if !c.crush {
					t.Fatalf("Error candy found %v but should not", c)
				}
			} else {
				if ok {
					assertCandyType(t, c._type, ct)
				} else {
					t.Fatalf("Error candy not found at %d,%d", curx, cury)
				}
			}
			curx += BlockSize
		}
		cury += BlockSize
		curx = XMin
	}
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

func assertNormal(t *testing.T, c *Candy, expected bool) {

	if c.IsNormal() != expected {
		t.Errorf("Wrong type %v, expected normal=%t but was %t", c, expected, c.IsNormal())
	}
}

func assertStriped(t *testing.T, c *Candy, expected bool) {
	if c.IsStriped() != expected {
		t.Errorf("Wrong type %v, expected striped=%t but was %t", c, expected, c.IsStriped())
	}
}

func assertStripedH(t *testing.T, c *Candy, expected bool) {
	if c.IsStripedH() != expected {
		t.Errorf("Wrong type %v, expected striped h=%t but was %t", c, expected, c.IsStripedH())
	}
}

func assertStripedV(t *testing.T, c *Candy, expected bool) {
	if c.IsStripedV() != expected {
		t.Errorf("Wrong type %v, expected striped v=%t but was %t", c, expected, c.IsStripedV())
	}
}

func assertPacked(t *testing.T, c *Candy, expected bool) {
	if c.IsPacked() != expected {
		t.Errorf("Wrong type %v, expected packed=%t but was %t", c, expected, c.IsStriped())
	}
}

func crushThem(them ...int) {
	for _, i := range them {
		g.candys[i].crush = true
	}
}

func assertCandy(t *testing.T, c, expected *Candy) {
	if c != expected {
		t.Errorf("Wrong candy, expected %v but was %v", expected, c)
	}
}
