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
	match := false
	// remove selection
	for _, c := range g.candys {
		c.visitedColumn = false
		c.visitedLine = false
	}
	//fmt.Println("check lines")
	for _, c := range g.candys {
		lines := g.findInLine(c, c._type)
		match = checkRegion(lines, false) || match
	}
	//fmt.Println("check columns")
	for _, c := range g.candys {
		columns := g.findInColumn(c, c._type)
		match = checkRegion(columns, true) || match
	}
	return match
}

func (g *Game) findInColumn(c *Candy, t CandyType) Region {
	return findInColumn(g.candys, nil, c, t)
}

func findInColumn(all, region Region, c *Candy, t CandyType) Region {
	if c == nil || c.visitedColumn || c._type != t {
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

func checkRegion(region Region, vertical bool) bool {
	nbMatch := len(region)
	if nbMatch > 2 {
		//fmt.Printf("match region %v\n", region)
		for _, c := range region {
			if c.crush == EmptyCandy {
				// first time the candy receives crush vote
				c.crush = CrushCandy
			} else {
				// more than one time the candy receivees a crush vote
				// it will be transformed to a Packed Candy
				c.crush = packedCandy(c._type)
			}
		}
		// only special candy here
		if nbMatch == 4 {
			region[0].crush = stripesCandy(region[0]._type, vertical)
		}
		if nbMatch > 4 {
			region[0].crush = BombCandy
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
	return t1 == t2
}
