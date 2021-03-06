package main

import (
	"testing"
)

func TestClickTopLeft(t *testing.T) {
	setup()
	fillGame()

	g.Click(XMin+5, 5)

	if g.selected == nil {
		t.Fatal("No candy selected")
	}
	assertXY(t, g.selected, XMin, YMin)
}

func TestClickMiddle(t *testing.T) {
	setup()
	fillGame()

	g.Click(XMin+3*BlockSize+5, BlockSize*4+5)

	if g.selected == nil {
		t.Fatal("No candy selected")
	}
	assertXY(t, g.selected, XMin+3*BlockSize, BlockSize*4)
}

func TestClickOutside(t *testing.T) {
	setup()
	fillGame()

	g.Click(5, 5)

	if g.selected != nil {
		t.Fatal("Click outside shouldn't select a candy")
	}
}

func TestClickNear(t *testing.T) {
	c := &Candy{x: XMin + 3*BlockSize, y: 3 * BlockSize}
	near1 := &Candy{x: XMin + 4*BlockSize, y: 3 * BlockSize}
	near2 := &Candy{x: XMin + 3*BlockSize, y: 2 * BlockSize}
	notNear1 := &Candy{x: XMin + 3*BlockSize, y: 3 * BlockSize}
	notNear2 := &Candy{x: XMin + 2*BlockSize, y: 2 * BlockSize}

	assertNear(t, near(c, near1), true, c, near1)
	assertNear(t, near(c, near1), true, c, near2)
	assertNear(t, near(c, notNear1), false, c, notNear1)
	assertNear(t, near(c, notNear2), false, c, notNear2)
}

func TestCollision(t *testing.T) {
	c1 := &Candy{x: 10, y: 20}
	c2 := &Candy{x: 5, y: 10}

	collision := collide(c1, c2)

	if !collision {
		t.Errorf("(%d,%d) and (%d,%d) should collide", c1.x, c1.y, c2.x, c2.y)
	}
}

func TestCollisionColumn(t *testing.T) {
	setup()
	g.candys = popCandys([][]CandyType{{RedPackedCandy, RedPackedCandy}})

	collision := g.collideColumn(g.candys[0])

	if collision {
		t.Errorf("(0,0) and (0,%d) should not collide", BlockSize)
	}
}

func TestNoCollision(t *testing.T) {
	c1 := &Candy{x: 100, y: 20}
	c2 := &Candy{x: 5, y: 10}

	collision := collide(c1, c2)

	if collision {
		t.Errorf("(%d,%d) and (%d,%d) should not collide", c1.x, c1.y, c2.x, c2.y)
	}

}

func TestPopulateDropZone(t *testing.T) {
	setup()

	g.populateDropZone()

	assertNbCandy(t, NbBlockWidth)
	assertNotEmpty(t, g.candys[0])
}

func TestGenerateCandy(t *testing.T) {
	setup()

	candy := g.newCandy()

	assertNotEmpty(t, candy)
}

func TestIsNormal(t *testing.T) {
	assertNormal(t, &Candy{_type: EmptyCandy}, false)

	assertNormal(t, &Candy{_type: RedCandy}, true)
	assertNormal(t, &Candy{_type: GreenCandy}, true)
	assertNormal(t, &Candy{_type: YellowCandy}, true)
	assertNormal(t, &Candy{_type: PinkCandy}, true)
	assertNormal(t, &Candy{_type: OrangeCandy}, true)
	assertNormal(t, &Candy{_type: BlueCandy}, true)

	assertNormal(t, &Candy{_type: RedHStripesCandy}, false)
	assertNormal(t, &Candy{_type: PinkPackedCandy}, false)
	assertNormal(t, &Candy{_type: BombCandy}, false)
}

func TestIsStriped(t *testing.T) {
	assertStriped(t, &Candy{_type: RedHStripesCandy}, true)
	assertStriped(t, &Candy{_type: RedVStripesCandy}, true)
	assertStriped(t, &Candy{_type: RedCandy}, false)
	assertStriped(t, &Candy{_type: RedPackedCandy}, false)
	assertStriped(t, &Candy{_type: BombCandy}, false)

	assertStripedH(t, &Candy{_type: RedHStripesCandy}, true)
	assertStripedH(t, &Candy{_type: RedVStripesCandy}, false)
	assertStripedH(t, &Candy{_type: RedCandy}, false)
	assertStripedH(t, &Candy{_type: RedPackedCandy}, false)
	assertStripedH(t, &Candy{_type: BombCandy}, false)

	assertStripedV(t, &Candy{_type: RedHStripesCandy}, false)
	assertStripedV(t, &Candy{_type: RedVStripesCandy}, true)
	assertStripedV(t, &Candy{_type: RedCandy}, false)
	assertStripedV(t, &Candy{_type: RedPackedCandy}, false)
	assertStripedV(t, &Candy{_type: BombCandy}, false)
}

func TestIsPacked(t *testing.T) {
	assertPacked(t, &Candy{_type: RedPackedCandy}, true)
	assertPacked(t, &Candy{_type: OrangePackedCandy}, true)
	assertPacked(t, &Candy{_type: PinkPackedCandy}, true)
	assertPacked(t, &Candy{_type: YellowPackedCandy}, true)
	assertPacked(t, &Candy{_type: BluePackedCandy}, true)
	assertPacked(t, &Candy{_type: GreenPackedCandy}, true)

	assertPacked(t, &Candy{_type: RedHStripesCandy}, false)
	assertPacked(t, &Candy{_type: BombCandy}, false)
	assertPacked(t, &Candy{_type: PinkVStripesCandy}, false)
	assertPacked(t, &Candy{_type: YellowCandy}, false)
}

func TestTopCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy},
		{RedCandy},
	})

	c := topCandy(candys, candys[1])

	assertCandy(t, c, candys[0])
}

func TestBottomCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy},
		{RedCandy},
	})

	c := bottomCandy(candys, candys[0])

	assertCandy(t, c, candys[1])
}

func TestLeftCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy, RedCandy},
	})

	c := leftCandy(candys, candys[1])

	assertCandy(t, c, candys[0])
}

func TestRightCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy, RedCandy},
	})

	c := rightCandy(candys, candys[0])

	assertCandy(t, c, candys[1])
}

func TestTopLeftCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy, RedCandy},
		{RedCandy, RedCandy},
	})

	c := topLeftCandy(candys, candys[3])

	assertCandy(t, c, candys[0])
}

func TestTopRightCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy, RedCandy},
		{RedCandy, RedCandy},
	})

	c := topRightCandy(candys, candys[2])

	assertCandy(t, c, candys[1])
}

func TestBottomLeftCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy, RedCandy},
		{RedCandy, RedCandy},
	})

	c := bottomLeftCandy(candys, candys[1])

	assertCandy(t, c, candys[2])
}

func TestBottomRightCandy(t *testing.T) {
	candys := popCandys([][]CandyType{
		{RedCandy, RedCandy},
		{RedCandy, RedCandy},
	})

	c := bottomRightCandy(candys, candys[0])

	assertCandy(t, c, candys[3])
}
