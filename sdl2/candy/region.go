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

const (
	Left, Top, Right, Bottom string = "0", "1", "2", "3"
)

//func regionSignature(region Region) string {
//	for _, c := range region {
//		c.visited = false
//	}
//	return candySignature(region, region[0])
//}
//
//func candySignature(region Region, c *Candy) string {
//	c.visited = true
//	var signature string
//	// left, top, right, bottom
//	tl := leftCandy(region, c)
//	if tl != nil && !tl.visited {
//		signature += Left + candySignature(region, tl)
//	}
//	tc := topCandy(region, c)
//	if tc != nil && !tc.visited {
//		signature += Top + candySignature(region, tc)
//	}
//	tr := rightCandy(region, c)
//	if tr != nil && !tr.visited {
//		signature += Right + candySignature(region, tr)
//	}
//	tb := bottomCandy(region, c)
//	if tb != nil && !tb.visited {
//		signature += Bottom + candySignature(region, tb)
//	}
//	return signature
//}
//
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

func crushable(region Region) bool {
	nb := len(region)
	if nb < 3 {
		return false
	}
	if nb == 3 && !alligned(region) {
		return false
	}
	return true
}
