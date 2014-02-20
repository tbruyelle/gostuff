package main

import (
	"testing"
)

func TestAlligned(t *testing.T) {
	candysXAlligned := generateCandys(C{d: Right}, C{d: Right}, C{d: Right})
	candysYAlligned := generateCandys(C{d: Bottom}, C{d: Bottom}, C{d: Bottom})
	candysNotAlligned := generateCandys(C{d: Bottom}, C{d: Left}, C{d: Bottom})

	if !alligned(candysXAlligned) {
		t.Errorf("Should be X alligned %+v", candysXAlligned)
	}
	if !alligned(candysYAlligned) {
		t.Errorf("Should be Y alligned %+v", candysYAlligned)
	}
	if alligned(candysNotAlligned) {
		t.Errorf("Should not be alligned %v", candysNotAlligned)
	}
}

func TestPermuteLeftRight(t *testing.T) {
	setup()
	fillGame()
	c1 := g.candys[0]
	c2 := g.candys[1]

	g.permute(c1, c2)

	assertVx(t, c1, BlockSize)
	assertVx(t, c2, -BlockSize)
}

func TestTranslateLeftRight(t *testing.T) {
	setup()
	fillGame()
	c1 := g.candys[0]
	ox1 := c1.x
	oy1 := c1.y
	c2 := g.candys[1]
	ox2 := c2.x
	oy2 := c2.y
	g.permute(c1, c2)

	g.translate()

	assertXY(t, c1, ox1+tSpeed, oy1)
	assertXY(t, c2, ox2-tSpeed, oy2)
}

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
	g.candys = generateCandys(C{}, C{})

	collision := g.collideColumn(g.candys[0], 0)

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

func TestFall(t *testing.T) {
	setup()
	g.populateDropZone()
	g.applyGravity()

	moving := g.fall()

	assertNbCandy(t, NbBlockWidth)
	assertY(t, g.candys[0], 1)
	if !moving {
		t.Error("Wrong move state, should still moving")
	}
}

func TestFalls(t *testing.T) {
	setup()
	g.populateDropZone()
	g.applyGravity()

	for g.fall() {
	}

	assertY(t, g.candys[0], WindowHeight-BlockSize)
}

func TestFallAll(t *testing.T) {
	setup()

	fillGame()

	// disable because the order is not respected because the gravity speed
	// is different threw columns
	//for i, c := range g.candys {
	//	x := XMin + BlockSize*(i%NbBlockWidth)
	//	fmt.Println(i, i%NbBlockWidth, x)
	//	y := WindowHeight - BlockSize*(1+(i/NbBlockHeight))
	//	assertXY(t, c, x, y)
	//}
	assertNbCandy(t, NbBlockWidth*NbBlockHeight)
}

func TestGenerateCandy(t *testing.T) {
	setup()

	candy := g.newCandy()

	assertNotEmpty(t, candy)
}

func TestApplyGravity(t *testing.T) {
	setup()
	g.populateDropZone()

	g.applyGravity()

	for i, c := range g.candys {
		assertGravity(t, c, i%2+1)
	}
}
