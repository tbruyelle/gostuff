package main

import "fmt"
import "testing"

func TestFindRegion(t *testing.T) {
	all := generateRegion(
		C{Right, RedCandy}, C{Right, RedCandy}, C{Right, RedCandy},
		C{Bottom, RedCandy}, C{Left, YellowCandy}, C{Left, YellowCandy},
		C{Bottom, BlueCandy}, C{Right, BlueCandy}, C{Right, RedCandy},
	)

	fmt.Printf("line=%v\n", findInLine(all, nil, all[0], all[0]._type))
	fmt.Printf("column=%v\n", findInColumn(all, nil, all[2], all[2]._type))
	//for _,c:= range all{
	//		region:=findInLine(all,nil,c,c._type)
	//	}
}

type C struct {
	dir string
	t   CandyType
}

func generateRegion(cs ...C) Region {
	region := []*Candy{&Candy{x: XMin, y: YMin}}
	curx, cury := XMin, YMin
	for _, c := range cs {
		switch c.dir {
		case Left:
			curx -= BlockSize
		case Right:
			curx += BlockSize
		case Top:
			cury -= BlockSize
		case Bottom:
			cury += BlockSize
		}
		region = append(region, &Candy{x: curx, y: cury, _type: c.t})
	}
	return region
}
