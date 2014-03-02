package main

// remove crushed candys
func (g *Game) crushing() {
	// first loop, handled crushed special candys
	processed := make(map[*Candy]bool)
	for g.crushSpecials(processed) {
		// recall crushSpecials() until there
		// is no more special candys to crush
	}
}

// crushSpecials() returns true when a special
// crush has been crushed, saying the loop needs to
// be restarted from the beginning.
// Indeed, a special candy crush may crushed other special
// candys that needs to be handled by that method
func (g *Game) crushSpecials(processed map[*Candy]bool) bool {
	for _, c := range g.candys {
		if _, done := processed[c]; !done && c.WillDie() && c._type > NbCandyType {
			// the candy needs a special crush
			if c.IsStriped() {
				g.crushStripes(c)
				processed[c] = true
				return true
			} else if c._type == BombCandy {
				g.crushBomb(c)
				processed[c] = true
				return true
			} else if c.IsPacked() {
				g.crushPacked(c)
				processed[c] = true
				return true
			}
		}
	}
	return false
}

// crushStripes() crushes candys in a line or
// column, according the stripes direction.
func (g *Game) crushStripes(c *Candy) {
	c.ChangeState(NewDyingState())
	if c.IsStripedH() {
		crushDir(g.candys, c, Left)
		crushDir(g.candys, c, Right)
	}
	if c.IsStripedV() {
		crushDir(g.candys, c, Top)
		crushDir(g.candys, c, Bottom)
	}
}

func crushDir(cs []*Candy, c *Candy, dir Direction) {
	cc := c
	i := 0
	for {
		i += 3
		cc = findCandyInDir(cs, cc, dir)
		if cc == nil {
			break
		}
		cc.ChangeState(NewDyingStateDelayed(i))
	}
}

// crushBomb() crushes all candys of a determined
// type.
func (g *Game) crushBomb(bomb *Candy) {
	if g.translation != nil {
		if g.translation.c1 == bomb {
			g.crushBombWith(bomb, g.translation.c2._type)
			return
		} else if g.translation.c2 == bomb {
			g.crushBombWith(bomb, g.translation.c1._type)
			return
		}
	}
	// The bomb is crushed because of a combo
	// we need to pick a random CandyType.
	g.crushBombWith(bomb, g.candyTypeGen.NewCandyType())
}

func (g *Game) crushBombWith(bomb *Candy, _type CandyType) {
	if _type == BombCandy {
		// both candys are Bombs, we remove everything
		for i, c := range g.candys {
			c.ChangeState(NewDyingStateDelayed(i))
		}
		return
	}
	// remove all candys with same type
	for i, c := range g.candys {
		if matchType(c._type, _type) {
			c.ChangeState(NewDyingStateDelayed(i))
		}
	}
}

// crushPacked() crushed candys that surround the
// packed candy.
func (g *Game) crushPacked(packed *Candy) {
	for dir := Direction(0); dir < NbDirections; dir++ {
		c := findCandyInDir(g.candys, packed, dir)
		if c != nil {
			c.ChangeState(NewDyingState())
		}
	}
}
