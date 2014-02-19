package main

import (
	"fmt"
)

type Region []*Candy

func (r Region) String() string {
	var s string
	for _, c := range r {
		s += fmt.Sprintf("(%d,%d)", c.x, c.y)
	}
	return s
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
	if c == nil || c.visitedLine || c._type != t {
		return region
	}
	c.visitedLine = true
	region = append(region, c)
	region = findInLine(all, region, leftCandy(all, c), t)
	region = findInLine(all, region, rightCandy(all, c), t)
	return region
}

func checkRegion(region Region) bool {
	nbMatch := len(region)
	if nbMatch > 2 {
		//fmt.Printf("match region %v\n", region)
		for _, c := range region {
			c.matching++
		}
		// only special candy here
		if nbMatch == 4 {
			region[0].matching = FourInLine
		}
		if nbMatch > 4 {
			region[0].matching = FiveInLine
		}
		return true
	}
	return false
}
