package main

import "testing"

func TestCrushingStripesH(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{BlueCandy, PinkCandy, OrangeCandy},
		{RedCandy, PinkHStripesCandy, RedCandy},
		{RedCandy, PinkCandy, RedCandy},
	})
	crushThem(1, 4, 7)

	g.crushing()

	assertNbCandy(t, 4)
	assertCandyTypes(t, [][]CandyType{
		{BlueCandy, EmptyCandy, OrangeCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{RedCandy, EmptyCandy, RedCandy},
	})
}

func TestCrushingStripesV(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{BlueCandy, PinkCandy, OrangeCandy},
		{RedCandy, RedVStripesCandy, RedCandy},
		{RedCandy, PinkCandy, RedCandy},
	})
	crushThem(3, 4, 5)

	g.crushing()

	assertNbCandy(t, 4)
	assertCandyTypes(t, [][]CandyType{
		{BlueCandy, EmptyCandy, OrangeCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{RedCandy, EmptyCandy, RedCandy},
	})
}

func TestCrushingStripesH2(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{PinkCandy, PinkHStripesCandy, PinkCandy},
	})
	crushThem(0, 1, 2)

	g.crushing()

	assertNbCandy(t, 0)
}

func TestCrushingBomb(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, BombCandy, BlueCandy, RedCandy, RedCandy},
	})
	g.translation = &Translation{g.candys[1], g.candys[0]}
	crushThem(0, 1)

	g.crushing()

	assertNbCandy(t, 1)
	assertCandyTypes(t, [][]CandyType{
		{EmptyCandy, EmptyCandy, BlueCandy, EmptyCandy, EmptyCandy},
	})
}

func TestCrushingBomb_permuteOtherSide(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, BombCandy, BlueCandy, RedCandy, RedCandy},
	})
	crushThem(0, 1)
	g.translation = &Translation{g.candys[0], g.candys[1]}

	g.crushing()

	assertNbCandy(t, 1)
	assertCandyTypes(t, [][]CandyType{
		{EmptyCandy, EmptyCandy, BlueCandy, EmptyCandy, EmptyCandy},
	})
}

func TestCrushingBombOnStripes(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{BlueCandy, RedCandy, BlueCandy},
		{YellowCandy, BombCandy, YellowCandy},
		{GreenCandy, BlueCandy, GreenCandy},
		{PinkCandy, RedCandy, RedCandy},
		{PinkCandy, RedHStripesCandy, OrangeCandy},
	})
	crushThem(1, 4)
	g.translation = &Translation{g.candys[1], g.candys[4]}

	g.crushing()

	assertNbCandy(t, 8)
	assertCandyTypes(t, [][]CandyType{
		{BlueCandy, EmptyCandy, BlueCandy},
		{YellowCandy, EmptyCandy, YellowCandy},
		{GreenCandy, BlueCandy, GreenCandy},
		{PinkCandy, EmptyCandy, EmptyCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
	})
}

func TestCrushingStripesOnBomb(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{BlueCandy, RedCandy, BlueCandy},
		{YellowCandy, BombCandy, YellowCandy},
		{GreenCandy, BlueCandy, GreenCandy},
		{PinkCandy, RedCandy, RedCandy},
		{PinkCandy, PinkVStripesCandy, PinkCandy},
	})
	crushThem(12, 13, 14)
	// redefine the newCandy func so we can define
	// the candy type chosen by the bomb match

	g.crushing()

	assertNbCandy(t, 6)
	assertCandyTypes(t, [][]CandyType{
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{YellowCandy, EmptyCandy, YellowCandy},
		{GreenCandy, EmptyCandy, GreenCandy},
		{PinkCandy, EmptyCandy, RedCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
	})
}
