package main

import "testing"

func TestMatchingNothing(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, BlueCandy}, C{Right, RedCandy})

	match := g.matching()

	if match {
		t.Fatalf("Should not find a match")
	}
	assertCrushes(t, g.candys, EmptyCandy)
}

func TestMatchingThreeInLine(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys, CrushCandy)
}

func TestMatchingThreeInLineColumn(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Bottom, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys, CrushCandy)
}

func TestMatchingFourInLine(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Bottom, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], RedHStripesCandy)
	assertCrushes(t, g.candys[1:], CrushCandy)
}

func TestMatchingPacked(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys[:1], CrushCandy)
	assertCrush(t, g.candys[2], RedPackedCandy)
	assertCrushes(t, g.candys[3:], CrushCandy)
}

func TestMatchingBomb(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], BombCandy)
	assertCrushes(t, g.candys[1:], CrushCandy)
}

func TestVStripesCandy(t *testing.T) {
	assertCandyType(t, stripesCandy(RedCandy, false), RedVStripesCandy)
	assertCandyType(t, stripesCandy(BlueCandy, false), BlueVStripesCandy)
	assertCandyType(t, stripesCandy(GreenCandy, false), GreenVStripesCandy)
	assertCandyType(t, stripesCandy(YellowCandy, false), YellowVStripesCandy)
	assertCandyType(t, stripesCandy(OrangeCandy, false), OrangeVStripesCandy)
	assertCandyType(t, stripesCandy(PinkCandy, false), PinkVStripesCandy)
}

func TestHStripesCandy(t *testing.T) {
	assertCandyType(t, stripesCandy(RedCandy, true), RedHStripesCandy)
	assertCandyType(t, stripesCandy(BlueCandy, true), BlueHStripesCandy)
	assertCandyType(t, stripesCandy(GreenCandy, true), GreenHStripesCandy)
	assertCandyType(t, stripesCandy(YellowCandy, true), YellowHStripesCandy)
	assertCandyType(t, stripesCandy(OrangeCandy, true), OrangeHStripesCandy)
	assertCandyType(t, stripesCandy(PinkCandy, true), PinkHStripesCandy)
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
