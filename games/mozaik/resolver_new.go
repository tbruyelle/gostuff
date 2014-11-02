package main

import (
	"container/heap"
	"fmt"
	"sync"
)

type Nodes []*Node

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Less(i, j int) bool {
	return ns[i].priority < ns[j].priority
}

func (ns Nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns *Nodes) Push(x interface{}) {
	node := x.(*Node)
	*ns = append(*ns, node)
}

func (ns *Nodes) Pop() interface{} {
	old := *ns
	n := len(old)
	node := old[n-1]
	*ns = old[0 : n-1]
	return node
}

type Node struct {
	// current switch
	s      int
	depth  int
	parent *Node
	// lvl represents a copy of the level
	// at the current node of the tree
	lvl      Level
	priority int
}

var signs map[string]bool
var mutex = new(sync.Mutex)

func addSign(sign string) {
	mutex.Lock()
	defer mutex.Unlock()
	signs[sign] = true
}

func hasSign(sign string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	return signs[sign]
}

func (n *Node) String() string {
	//return fmt.Sprintf("s%d, d=%d, childs=%+v", n.s, n.depth, n.childs)
	//return fmt.Sprintf("s%d, d=%d, parent=[%+v] win=%t", n.s, n.depth, n.parent, n.lvl.Win())
	depth := n.depth
	return fmt.Sprintf("d=%d, p=%d, sws=%s", depth, n.priority, n.road())
}

// Returns the switch combination used so far
func (n *Node) road() string {
	var s string
	for n.parent != nil && n.s >= 0 {
		s = n.lvl.switches[n.s].name + s
		n = n.parent
	}
	if n.s >= 0 {
		s = n.lvl.switches[n.s].name + s
	}
	return s
}

func Resolve(lvl Level) *Node {
	ns := make(Nodes, 0)
	heap.Init(&ns)

	init := &Node{
		s:        -1,
		depth:    0,
		lvl:      lvl.Copy(),
		priority: lvl.HowFar(),
	}
	heap.Push(&ns, init)

	for {
		n := process(&ns)
		if n != nil {
			return n
		}
	}
	return nil
}

func process(ns *Nodes) *Node {
	n := heap.Pop(ns).(*Node)
	fmt.Println("Processing node", n)
	if n.lvl.Win() {
		return n
	}
	if n.depth > 50 {
		return nil
	}
	for i, _ := range n.lvl.switches {
		if n.lvl.IsPlain(i) {
			// Ignore switch with plain color
			continue
		}
		nn := &Node{
			s:        i,
			depth:    n.depth + 1,
			lvl:      n.lvl.Copy(),
			parent:   n,
			priority: n.lvl.HowFar() + n.depth,
		}
		nn.lvl.RotateSwitch(n.lvl.switches[i])
		heap.Push(ns, nn)
		fmt.Println("Added node", nn)
	}
	return nil
}
