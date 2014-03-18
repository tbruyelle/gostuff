package main

import (
	"math"
)

type Node struct {
	// current switch
	s      int
	depth  int
	childs []*Node
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
