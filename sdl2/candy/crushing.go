package main

import (
	"fmt"
)

// remove crushed candys
func (g *Game) crushing() {
	// first loop, handled crushed special candys
	for _, c := range g.candys {
		if c.crush && c._type > NbCandyType {
			// the candy needs a special crush
			if c.isStriped() {
				g.crushStripes(c)
			} else if c._type == BombCandy {
				g.crushBomb(c)
			}
		}
	}
	// final loop, remove the crushed candys
	var kept []*Candy
	for _, c := range g.candys {
		//fmt.Printf("crushCandy %v\n", c)
		if !c.crush {
			kept = append(kept, c)
		}
	}
	fmt.Printf("Crushing %d candys\n", len(g.candys)-len(kept))
	g.candys = kept
	//fmt.Printf("NOW %d candys\n", len(g.candys))
}

func (g *Game) crushStripes(c *Candy) {
	c.crush = true
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
	}
}

func (g *Game) crushBomb(bomb *Candy) {
	if g.translation != nil {
		if g.translation.c1._type == BombCandy {
			g.crushBombWith(bomb, g.translation.c2._type)
		} else if g.translation.c2._type == BombCandy {
			g.crushBombWith(bomb, g.translation.c1._type)
		}else {
			// the bomb is crushed because of a combo
			// we need to pick a random candytype
			//TODO
		}
	}
}

func (g *Game) crushBombWith(bomb *Candy, _type CandyType) {
	bomb.crush = true
	// remove all candys with same type
	for _, c := range g.candys {
		if matchType(c._type, _type) {
			c.crush = true
		}
	}
}
