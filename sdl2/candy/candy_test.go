package main

import (
	"testing"
)

var g *Game

func setup() {
	g = NewGame()
}

func TestGenerateCandy(t *testing.T) {
	setup()

	candy := g.NewCandy()

	if candy == EmptyCandy {
		t.Errorf("Wrong candy, expected not empty but was empty")
	}
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

func assertMatch(t *testing.T, match Match, start, length int) {
	if match.start != start && match.length != length {
		t.Errorf("Wrong match, expected start=%d, length=%d, but was (%d,%d)", start, length, match.start, match.length)
	}
}
