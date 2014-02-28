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

	assertCandyTypes(t, [][]CandyType{
		{EmptyCandy, EmptyCandy, EmptyCandy},
	})
}

func TestCrushingBomb(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{RedCandy, BombCandy, BlueCandy, RedCandy, RedCandy},
	})
	g.translation = &Translation{g.candys[1], g.candys[0]}
	crushThem(0, 1)

	g.crushing()

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

	assertCandyTypes(t, [][]CandyType{
		{BlueCandy, EmptyCandy, BlueCandy},
		{YellowCandy, EmptyCandy, YellowCandy},
		{GreenCandy, BlueCandy, GreenCandy},
		{PinkCandy, EmptyCandy, EmptyCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
	})
}

// Tests that a stripes which crushs a BombCandy makes
// the bomb explodes with a randomly chosen CandyType.
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
	// Redefine the CandyTypeGenerator with a mock
	// which will returns only BlueCandy.
	// So the BlueCandy will be the type chosen by
	// the BombCandy.
	g.candyTypeGen = candyTypeGenMock{BlueCandy}

	g.crushing()

	assertCandyTypes(t, [][]CandyType{
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{YellowCandy, EmptyCandy, YellowCandy},
		{GreenCandy, EmptyCandy, GreenCandy},
		{PinkCandy, EmptyCandy, RedCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
	})
}

func TestCrushingPacked(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{BlueCandy, RedCandy, BlueCandy},
		{YellowCandy, RedPackedCandy, YellowCandy},
		{PinkCandy, RedCandy, PinkCandy},
	})
	crushThem(4)

	g.crushing()

	assertCandyTypes(t, [][]CandyType{
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
	})
}

func TestCrushBombOnBomb(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{
		{BombCandy, BombCandy, RedCandy},
		{GreenCandy, YellowCandy, BlueCandy},
		{PinkCandy, RedCandy, PinkCandy},
	})
	crushThem(0, 1)
	g.translation = &Translation{g.candys[0], g.candys[1]}

	g.crushing()

	assertCandyTypes(t, [][]CandyType{
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
		{EmptyCandy, EmptyCandy, EmptyCandy},
	})
}
