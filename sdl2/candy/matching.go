package main

import (
	"fmt"
)

type Region []*Candy

func (r Region) String() string {
	var s string
	for _, c := range r {
		s += fmt.Sprintf("[%v]", c)
	}
	return s
}

func (g *Game) matching() bool {
	if g.translateBomb() {
		return true
	}
	match := false
	// remove selection
	for _, c := range g.candys {
		c.visitedColumn = false
		c.visitedLine = false
	}
	//fmt.Println("check lines")
	for _, c := range g.candys {
		lines := g.findInLine(c, c._type)
		match = g.checkRegion(lines, false) || match
	}
	//fmt.Println("check columns")
	for _, c := range g.candys {
		columns := g.findInColumn(c, c._type)
		match = g.checkRegion(columns, true) || match
	}

	return match
}

func (g *Game) findInColumn(c *Candy, t CandyType) Region {
	return findInColumn(g.candys, nil, c, t)
}

func findInColumn(all, region Region, c *Candy, t CandyType) Region {
	if c == nil || c.visitedColumn || !matchType(c._type, t) {
		return region
	}
	c.visitedColumn = true
	region = append(region, c)
	region = findInColumn(all, region, topCandy(all, c), t)
	region = findInColumn(all, region, bottomCandy(all, c), t)
	return region
}

func (g *Game) findInLine(c *Candy, t CandyType) Region {
	return findInLine(g.candys, nil, c, t)
}

func findInLine(all, region Region, c *Candy, t CandyType) Region {
	if c == nil || c.visitedLine || !matchType(c._type, t) {
		return region
	}
	c.visitedLine = true
	region = append(region, c)
	region = findInLine(all, region, leftCandy(all, c), t)
	region = findInLine(all, region, rightCandy(all, c), t)
	return region
}

// checkRegion() checks the region can produce a match.
// If yes it verifies also if the match can generate special
// candys.
func (g *Game) checkRegion(region Region, vertical bool) bool {
	nbMatch := len(region)
	if nbMatch > 2 {
		//fmt.Printf("match region %v\n", region)
		for _, c := range region {
			if !c.crush {
				// first time the candy receives crush vote
				c.crush = true
			} else if c.isNormal() {
				// more than one time the candy receivees a crush vote
				// it will be transformed to a Packed Candy
				c._type = packedCandy(c._type)
				c.crush = false
			}
		}
		if nbMatch >= 4 {
			// Handle special candy here, according
			// the match, we will mutate a simple
			// candy to a special one
			c := g.determineMutableCandy(region)
			// The mutation will occur only if we were able to
			// determinate a mutable candy.
			if c != nil {
				if nbMatch == 4 {
					// mutate candy to Stripes
					c._type = stripesCandy(region[0]._type, vertical)
					c.crush = false
				}
				if nbMatch > 4 {
					// mutate candy to Bomb
					c._type = BombCandy
					c.crush = false
				}
			}
		}
		return true
	}
	return false
}

// determineCandyToMutate() determines the candy in region
// which will mutate according to the special match triggered
func (g *Game) determineMutableCandy(region Region) *Candy {
	c, found := g.findTranslated(region)
	if !found {
		//find the first normal candy in the region
		for i := 0; i < len(region); i++ {
			if region[i].isNormal(){
				c = region[i]
				break
			}
		}
	}
	return c
}

func (g *Game) findTranslated(region Region) (*Candy, bool) {
	if g.translation != nil {
		if isInRegion(region, g.translation.c1) {
			return g.translation.c1, true
		}
		if isInRegion(region, g.translation.c2) {
			return g.translation.c2, true
		}
	}
	return nil, false
}

func isInRegion(region Region, c *Candy) bool {
	for i := 0; i < len(region); i++ {
		if region[i] == c {
			return true
		}
	}
	return false
}

func stripesCandy(t CandyType, vertical bool) CandyType {
	if vertical {
		return t + NbCandyType*2
	}
	return t + NbCandyType
}

func packedCandy(t CandyType) CandyType {
	return t + NbCandyType*3
}

func matchType(t1, t2 CandyType) bool {
	if t1 > NbCandyType*4 || t2 > NbCandyType*4 {
		return false
	}
	// compare type % nbCandyType to match stripes and packed candys
	return t1%NbCandyType == t2%NbCandyType
}

func (g *Game) translateBomb() bool {
	if g.translation != nil {
		if g.translation.c1._type == BombCandy || g.translation.c2._type == BombCandy {
			g.translation.c1.crush = true
			g.translation.c2.crush = true
			return true
		}
	}
	return false
}
