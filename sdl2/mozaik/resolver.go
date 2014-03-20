package main

import (
	"fmt"
	"math"
)

const (
	MaxDepth = 30
)

type Node struct {
	// current switch
	s      int
	depth  int
	parent *Node
	childs []*Node
	// lvl represents a copy of the level
	// at the current node of the tree
	lvl Level
}

func (n *Node) String() string {
	//return fmt.Sprintf("s%d, d=%d, childs=%+v", n.s, n.depth, n.childs)
	return fmt.Sprintf("s%d, d=%d, parent=[%+v] win=%t", n.s, n.depth, n.parent, n.lvl.Win())
}

var nearSw map[int][]int

func FindPaths(lvl Level) []*Node {
	nearSw = DetermineNearSwitches(lvl)
	paths := make([]*Node, len(lvl.switches))
	for i := range lvl.switches {
		fmt.Printf("find path starting switch %d\n", i)
		n := &Node{s: i, depth: 1, lvl: lvl.Copy()}
		paths[i] = n
		check(n)
		//fmt.Printf("switch %d path=%+v\n", i, n)
	}
	return paths
}

func check(n *Node) {
	if n.depth > MaxDepth {
		return
	}
	if n.lvl.Win() {
		fmt.Printf("WIN %+v\n", n)
		return
	}
	depth := n.depth + 1
	n.lvl.RotateSwitch(n.lvl.switches[n.s])
	if !n.hasRotatedMoreThanThrice() {
		n.childs = append(n.childs, &Node{s: n.s, depth: depth, parent: n, lvl: n.lvl.Copy()})
	}
	for _, ns := range nearSw[n.s] {
		n.childs = append(n.childs, &Node{s: ns, depth: depth, parent: n, lvl: n.lvl.Copy()})
	}
	//fmt.Printf("Check %d childs depth=%d\n", len(n.childs), n.depth)
	for _, c := range n.childs {
		check(c)
	}
}

func (n *Node) hasRotatedMoreThanThrice() bool {
	if n.parent != nil && n.parent.s == n.s && n.parent.parent != nil && n.parent.parent.s == n.s {
		return true
	}
	return false
}

// DetermineNearSwitches returns a map of nearest switches
// for each switches in the lvl in parameter
func DetermineNearSwitches(lvl Level) map[int][]int {
	res := make(map[int][]int, len(lvl.switches))
	for i := range lvl.switches {
		for j := range lvl.switches {
			if i == j {
				continue
			}
			lineDiff := math.Abs(float64(lvl.switches[i].line - lvl.switches[j].line))
			colDiff := math.Abs(float64(lvl.switches[i].col - lvl.switches[j].col))
			if lineDiff < 2 && colDiff < 2 {
				res[i] = append(res[i], j)
			}
		}

	}
	return res
}
