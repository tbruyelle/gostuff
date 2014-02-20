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
	assertCrush(t, g.candys[0], RedStripesCandy)
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

