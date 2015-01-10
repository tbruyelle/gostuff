package main

import (
	"fmt"
)

type ColorDef int

const (
	Red         ColorDef = iota //0
	Yellow                      //1
	Blue                        //2
	Green                       //3
	Pink                        //4
	Orange                      //5
	LightBlue                   //6
	Purple                      //7
	Brown                       //8
	LightGreen                  //9
	Cyan                        //A
	LightPink                   //B
	White                       //C
	LightPurple                 //D
	LightBrown                  //E
	OtherWhite                  //F
	Empty       ColorDef = -1
)

type Block struct {
	Color    ColorDef
	Rendered bool
}

type Switch struct {
	state     State
	line, col int
	X, Y, Z   int
	rotate    int
	name      string
}

func (s *Switch) Rotate() {
	s.ChangeState(NewRotateState())
}

func (s *Switch) ChangeState(state State) {
	if s.state != nil {
		s.state.Exit(g, s)
		if !s.state.AllowChange(state) {
			fmt.Println("Change state not allowed")
			return
		}
	}
	s.state = state
	s.state.Enter(g, s)
}

func (s *Switch) String() string {
	return fmt.Sprintf("sw{line:%d, col:%d}", s.line, s.col)
}
