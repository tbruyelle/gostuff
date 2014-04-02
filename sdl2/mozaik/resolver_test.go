package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindPaths_Level1(t *testing.T) {
	lvl := LoadLevel(1)

	//paths := FindShortestPaths(lvl)
	n := FindPathAsync(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestFindPaths_Level2(t *testing.T) {
	lvl := LoadLevel(2)

	//paths := FindShortestPaths(lvl)
	n := FindPathAsync(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestFindPaths_Level3(t *testing.T) {
	lvl := LoadLevel(3)

	//paths := FindShortestPaths(lvl)
	n := FindPathAsync(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestFindPaths_Level4(t *testing.T) {
	lvl := LoadLevel(4)

	//paths := FindShortestPaths(lvl)
	n := FindPathAsync(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestFindPaths_Level5(t *testing.T) {
	lvl := LoadLevel(5)

	//paths := FindShortestPaths(lvl)
	n := FindPathAsync(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestFindPaths_Level6(t *testing.T) {
	lvl := LoadLevel(6)

	//paths := FindShortestPaths(lvl)
	n := FindPathAsync(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestFindPaths_Level7(t *testing.T) {
	lvl := LoadLevel(7)

	//paths := FindShortestPaths(lvl)
	n := FindPathAsync(lvl)

	fmt.Printf("test result %+v\n", n)
}

func TestDetermineNearestSwicthes_Level1(t *testing.T) {
	lvl := LoadLevel(1)

	res := DetermineNearSwitches(lvl)

	assert.Equal(t, len(lvl.switches), len(res))
	assert.Equal(t, len(res[0]), 1)
	assert.Equal(t, len(res[1]), 2)
	assert.Equal(t, len(res[2]), 1)
	assert.Equal(t, res[0], []int{1})
	assert.Equal(t, res[1], []int{0, 2})
	assert.Equal(t, res[2], []int{1})
}

func TestDetermineNearestSwicthes_Level3(t *testing.T) {
	lvl := LoadLevel(3)

	res := DetermineNearSwitches(lvl)

	assert.Equal(t, len(lvl.switches), len(res))
	assert.Equal(t, len(res[0]), 3)
	assert.Equal(t, len(res[1]), 5)
	assert.Equal(t, len(res[2]), 3)
	assert.Equal(t, len(res[3]), 5)
	assert.Equal(t, len(res[4]), 8)
	assert.Equal(t, len(res[5]), 5)
	assert.Equal(t, len(res[6]), 3)
	assert.Equal(t, len(res[7]), 5)
	assert.Equal(t, len(res[8]), 3)
	assert.Equal(t, res[0], []int{1, 3, 4})
	assert.Equal(t, res[1], []int{0, 2, 3, 4, 5})
	assert.Equal(t, res[2], []int{1, 4, 5})
	assert.Equal(t, res[3], []int{0, 1, 4, 6, 7})
	assert.Equal(t, res[4], []int{0, 1, 2, 3, 5, 6, 7, 8})
	assert.Equal(t, res[5], []int{1, 2, 4, 7, 8})
	assert.Equal(t, res[6], []int{3, 4, 7})
	assert.Equal(t, res[7], []int{3, 4, 5, 6, 8})
	assert.Equal(t, res[8], []int{4, 5, 7})
}
