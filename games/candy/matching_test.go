package main

import "testing"

func TestMatchingNothing(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, BlueCandy, RedCandy},
	})

	match := g.matching()

	if match {
		t.Fatalf("Should not find a match")
	}
	assertCrushes(t, g.candys, false, EmptyCandy)
}

func TestMatchingThreeInLine(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, RedCandy, RedCandy},
	})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys, true, RedCandy)
}

func TestMatchingThreeInColumn(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy},
		{RedCandy},
		{RedCandy},
	})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys, true, RedCandy)
}

func TestMatchingStripes(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, RedCandy, RedCandy, RedCandy},
	})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], false, RedHStripesCandy)
	assertCrushes(t, g.candys[1:], true, RedCandy)
}

func TestMatchingPacked(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, RedCandy, RedCandy},
		{BlueCandy, YellowCandy, RedCandy},
		{OrangeCandy, PinkCandy, RedCandy},
	})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys[:1], true, RedCandy)
	assertCrush(t, g.candys[2], false, RedPackedCandy)
	assertCrushes(t, g.candys[3:4], false, EmptyCandy)
	assertCrush(t, g.candys[5], true, RedCandy)
	assertCrushes(t, g.candys[6:7], false, EmptyCandy)
	assertCrush(t, g.candys[8], true, RedCandy)
}

func TestMatchingBomb(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, RedCandy, RedCandy, RedCandy, RedCandy},
	})

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], false, BombCandy)
	assertCrushes(t, g.candys[1:], true, RedCandy)
}

func TestVStripesCandy(t *testing.T) {
	assertCandyType(t, stripesCandy(RedCandy, true), RedVStripesCandy)
	assertCandyType(t, stripesCandy(BlueCandy, true), BlueVStripesCandy)
	assertCandyType(t, stripesCandy(GreenCandy, true), GreenVStripesCandy)
	assertCandyType(t, stripesCandy(YellowCandy, true), YellowVStripesCandy)
	assertCandyType(t, stripesCandy(OrangeCandy, true), OrangeVStripesCandy)
	assertCandyType(t, stripesCandy(PinkCandy, true), PinkVStripesCandy)
}

func TestHStripesCandy(t *testing.T) {
	assertCandyType(t, stripesCandy(RedCandy, false), RedHStripesCandy)
	assertCandyType(t, stripesCandy(BlueCandy, false), BlueHStripesCandy)
	assertCandyType(t, stripesCandy(GreenCandy, false), GreenHStripesCandy)
	assertCandyType(t, stripesCandy(YellowCandy, false), YellowHStripesCandy)
	assertCandyType(t, stripesCandy(OrangeCandy, false), OrangeHStripesCandy)
	assertCandyType(t, stripesCandy(PinkCandy, false), PinkHStripesCandy)
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

	assertMatch(t, matchType(RedCandy, BombCandy), false)
	assertMatch(t, matchType(GreenCandy, BombCandy), false)
	assertMatch(t, matchType(BlueCandy, BombCandy), false)
	assertMatch(t, matchType(OrangeCandy, BombCandy), false)
	assertMatch(t, matchType(YellowCandy, BombCandy), false)
	assertMatch(t, matchType(PinkCandy, BombCandy), false)
}

func TestMatchingWithBomb(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, BombCandy},
	})
	g.translation = &Translation{g.candys[0], g.candys[1]}

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrush(t, g.candys[0], true, RedCandy)
	assertCrush(t, g.candys[1], true, BombCandy)
}

func TestMatchingStripes_withTranslation(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, RedCandy, RedCandy, RedCandy},
		{BlueCandy, GreenCandy, OrangeCandy, PinkCandy},
	})
	g.translation = &Translation{g.candys[2], g.candys[6]}

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys[:1], true, RedCandy)
	assertCrush(t, g.candys[2], false, RedHStripesCandy)
	assertCrush(t, g.candys[3], true, RedCandy)
	assertCrushes(t, g.candys[4:], false, EmptyCandy)
}

func TestMatchingBomb_withTranslation(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, RedCandy, RedCandy, RedCandy, RedCandy},
		{BlueCandy, BlueCandy, GreenCandy, GreenCandy, YellowCandy},
	})
	g.translation = &Translation{g.candys[2], g.candys[7]}

	match := g.matching()

	if !match {
		t.Fatalf("Should find a match")
	}
	assertCrushes(t, g.candys[:1], true, RedCandy)
	assertCrush(t, g.candys[2], false, BombCandy)
	assertCrushes(t, g.candys[3:4], true, RedCandy)
	assertCrushes(t, g.candys[5:], false, EmptyCandy)
}
