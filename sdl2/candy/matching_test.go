package main

import "testing"

func TestMatchingNothing(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, BlueCandy}, C{Right, RedCandy})

	match := g.matching()

	if match {
		t.Fatalf("Should not find a match")
	}
	assertCrushes(t, g.candys, false, EmptyCandy)
}

func TestMatchingThreeInLine(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys, true, EmptyCandy)
}

func TestMatchingThreeInLineColumn(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Bottom, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys, true, EmptyCandy)
}

func TestMatchingFourInLine(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Bottom, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], true, RedHStripesCandy)
	assertCrushes(t, g.candys[1:], true, EmptyCandy)
}

func TestMatchingPacked(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys[:1], true, EmptyCandy)
	assertCrush(t, g.candys[2], true, RedPackedCandy)
	assertCrushes(t, g.candys[3:], true, EmptyCandy)
}

func TestMatchingBomb(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], true, BombCandy)
	assertCrushes(t, g.candys[1:], true, EmptyCandy)
}

func TestMatchWithStripesH(t *testing.T) {
	setup()
	g.candys = generateCandys(
		C{Right, BlueCandy}, C{Right, PinkCandy}, C{Right, OrangeCandy},
		C{Bottom, RedCandy}, C{Left, PinkHStripesCandy}, C{Left, RedCandy},
		C{Bottom, RedCandy}, C{Right, PinkCandy}, C{Right, RedCandy},
	)

	match := g.matching()
	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], false, EmptyCandy)
	assertCrush(t, g.candys[1], true, EmptyCandy)
	assertCrush(t, g.candys[2], false, EmptyCandy)
	assertCrushes(t, g.candys[3:5], true, UnmutableCandy)
	assertCrush(t, g.candys[6], false, EmptyCandy)
	assertCrush(t, g.candys[7], true, EmptyCandy)
	assertCrush(t, g.candys[8], false, EmptyCandy)
}

func TestMatchWithStripesV(t *testing.T) {
	setup()
	g.candys = generateCandys(
		C{Right, BlueCandy}, C{Right, PinkCandy}, C{Right, OrangeCandy},
		C{Bottom, RedCandy}, C{Left, RedVStripesCandy}, C{Left, RedCandy},
		C{Bottom, RedCandy}, C{Right, PinkCandy}, C{Right, RedCandy},
	)

	match := g.matching()
	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], false, EmptyCandy)
	assertCrush(t, g.candys[1], true, UnmutableCandy)
	assertCrush(t, g.candys[2], false, EmptyCandy)
	assertCrush(t, g.candys[3], true, EmptyCandy)
	assertCrush(t, g.candys[4], true, UnmutableCandy)
	assertCrush(t, g.candys[5], true, EmptyCandy)
	assertCrush(t, g.candys[6], false, EmptyCandy)
	assertCrush(t, g.candys[7], true, UnmutableCandy)
	assertCrush(t, g.candys[8], false, EmptyCandy)
}

func TestMatchWithStripesH2(t *testing.T) {
	setup()
	g.candys = generateCandys(
		C{Right, PinkCandy}, C{Right, PinkHStripesCandy}, C{Right, PinkCandy},
	)

	match := g.matching()
	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys, true, UnmutableCandy)
}

func TestMatchingWithBomb(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Bottom, RedCandy}, C{Bottom, BombCandy}, C{Bottom, BlueCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})
	g.translation = &Translation{g.candys[1], g.candys[0]}

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], true, UnmutableCandy)
	assertCrush(t, g.candys[1], true, EmptyCandy)
	assertCrushes(t, g.candys[3:], true, UnmutableCandy)
	assertCrush(t, g.candys[2], false, EmptyCandy)
}

func TestMatchingWithBomb_otherSide(t *testing.T) {
	setup()
	g.candys = generateCandys(C{Bottom, RedCandy}, C{Bottom, BombCandy}, C{Bottom, BlueCandy}, C{Bottom, RedCandy}, C{Bottom, RedCandy})
	g.translation = &Translation{g.candys[0], g.candys[1]}

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], true, UnmutableCandy)
	assertCrush(t, g.candys[1], true, EmptyCandy)
	assertCrush(t, g.candys[2], false, EmptyCandy)
	assertCrushes(t, g.candys[3:], true, UnmutableCandy)
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

	assertMatch(t, matchType(RedCandy, RedVStripesCandy), true)
	assertMatch(t, matchType(GreenCandy, GreenVStripesCandy), true)
	assertMatch(t, matchType(BlueCandy, BlueVStripesCandy), true)
	assertMatch(t, matchType(OrangeCandy, OrangeVStripesCandy), true)
	assertMatch(t, matchType(YellowCandy, YellowVStripesCandy), true)
	assertMatch(t, matchType(PinkCandy, PinkVStripesCandy), true)

	assertMatch(t, matchType(RedCandy, RedHStripesCandy), true)
	assertMatch(t, matchType(GreenCandy, GreenHStripesCandy), true)
	assertMatch(t, matchType(BlueCandy, BlueHStripesCandy), true)
	assertMatch(t, matchType(OrangeCandy, OrangeHStripesCandy), true)
	assertMatch(t, matchType(YellowCandy, YellowHStripesCandy), true)
	assertMatch(t, matchType(PinkCandy, PinkHStripesCandy), true)

	assertMatch(t, matchType(RedCandy, RedPackedCandy), true)
	assertMatch(t, matchType(GreenCandy, GreenPackedCandy), true)
	assertMatch(t, matchType(BlueCandy, BluePackedCandy), true)
	assertMatch(t, matchType(OrangeCandy, OrangePackedCandy), true)
	assertMatch(t, matchType(YellowCandy, YellowPackedCandy), true)
	assertMatch(t, matchType(PinkCandy, PinkPackedCandy), true)
}
