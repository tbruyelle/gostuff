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

func (g *Game) translateBomb() bool {
	if g.translation != nil {
		if g.translation.c1._type == BombCandy {
			g.crushBomb(g.translation.c1, g.translation.c2)
			return true
		}
		if g.translation.c2._type == BombCandy {
			g.crushBomb(g.translation.c2, g.translation.c1)
			return true
		}
	}
	return false
}

func (g *Game) crushBomb(bomb *Candy, other *Candy) {
	bomb.crush = true
	// remove all candys with same type
	for _, c := range g.candys {
		if matchType(c._type, other._type) {
			c.crush = true
			c.mutation = UnmutableCandy
		}
	}
}

func (g *Game) crushStripes(c *Candy) {
	c.crush = true
	c.mutation = UnmutableCandy
	if c.isStripedH() {
		crushDir(g.candys, c, Left)
		crushDir(g.candys, c, Right)
	}
	if c.isStripedV() {
		crushDir(g.candys, c, Top)
		crushDir(g.candys, c, Bottom)
	}
}

func crushDir(cs []*Candy, c *Candy, dir Direction) {
	cc := c
	for {
		cc = findCandyInDir(cs, cc, dir)
		if cc == nil {
			break
		}
		cc.crush = true
		cc.mutation = UnmutableCandy
	}
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

func (g *Game) checkRegion(region Region, vertical bool) bool {
	nbMatch := len(region)
	if nbMatch > 2 {
		//fmt.Printf("match region %v\n", region)
		for _, c := range region {
			if c.isStriped() {
				g.crushStripes(c)
			}
			if !c.crush {
				// first time the candy receives crush vote
				c.crush = true
			} else if c.mutation == EmptyCandy {
				// more than one time the candy receivees a crush vote
				// it will be transformed to a Packed Candy
				c.mutation = packedCandy(c._type)
			}
		}
		// only special candy here
		if nbMatch == 4 {
			region[0].mutation = stripesCandy(region[0]._type, vertical)
		}
		if nbMatch > 4 {
			region[0].mutation = BombCandy
		}
		return true
	}
	return false
}

func stripesCandy(t CandyType, vertical bool) CandyType {
	if vertical {
		return t + NbCandyType
	}
	return t + NbCandyType*2
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
