package main

import "testing"

func TestStripesCandy(t *testing.T) {
	assertCandyType(t, stripesCandy(RedCandy), RedStripesCandy)
	assertCandyType(t, stripesCandy(BlueCandy), BlueStripesCandy)
	assertCandyType(t, stripesCandy(GreenCandy), GreenStripesCandy)
	assertCandyType(t, stripesCandy(YellowCandy), YellowStripesCandy)
	assertCandyType(t, stripesCandy(OrangeCandy), OrangeStripesCandy)
	assertCandyType(t, stripesCandy(PinkCandy), PinkStripesCandy)
}

func TestPackedCandy(t *testing.T) {
	assertCandyType(t, packedCandy(RedCandy), RedPackedCandy)
	assertCandyType(t, packedCandy(BlueCandy), BluePackedCandy)
	assertCandyType(t, packedCandy(GreenCandy), GreenPackedCandy)
	assertCandyType(t, packedCandy(YellowCandy), YellowPackedCandy)
	assertCandyType(t, packedCandy(OrangeCandy), OrangePackedCandy)
	assertCandyType(t, packedCandy(PinkCandy), PinkPackedCandy)
}

func TestMatchType(t *testing.T) {
	assertMatch(t, matchType(RedCandy, RedCandy), true)
	assertMatch(t, matchType(GreenCandy, GreenCandy), true)
	assertMatch(t, matchType(BlueCandy, BlueCandy), true)
	assertMatch(t, matchType(OrangeCandy, OrangeCandy), true)
	assertMatch(t, matchType(YellowCandy, YellowCandy), true)
	assertMatch(t, matchType(PinkCandy, PinkCandy), true)

	assertMatch(t, matchType(RedCandy, PinkCandy), false)
	assertMatch(t, matchType(RedCandy, GreenCandy), false)
	assertMatch(t, matchType(RedCandy, OrangeCandy), false)
	assertMatch(t, matchType(RedCandy, YellowCandy), false)
	assertMatch(t, matchType(RedCandy, BlueCandy), false)
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
