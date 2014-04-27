package main

import "testing"

func TestReadVertexFile(t *testing.T) {

	vs := readVertexFile("data/cube")

	if len(vs) != 36 {
		t.Errorf("wrong number of vertices, expected 36 but have %d", len(vs))
	}
	//for _, v := range vs {
	//	fmt.Printf("%+v\n", v)
	//}
}

func TestSequence(t *testing.T) {

	assertSequence(t, 0, Sequence(3, 0), 0)
	assertSequence(t, 0, Sequence(3, 1), 1)
	assertSequence(t, 0, Sequence(3, 2), 2)

	assertSequence(t, 1, Sequence(3, 3), 3)
	assertSequence(t, 1, Sequence(3, 4), 4)
	assertSequence(t, 1, Sequence(3, 5), 5)

	assertSequence(t, 2, Sequence(3, 6), 6)
	assertSequence(t, 2, Sequence(3, 7), 7)
	assertSequence(t, 2, Sequence(3, 8), 8)

	assertSequence(t, 0, Sequence(3, 9), 9)
	assertSequence(t, 0, Sequence(3, 10), 10)
	assertSequence(t, 0, Sequence(3, 11), 11)

	assertSequence(t, 1, Sequence(3, 12), 12)
	assertSequence(t, 1, Sequence(3, 13), 14)
	assertSequence(t, 1, Sequence(3, 14), 14)

	assertSequence(t, 2, Sequence(3, 15), 15)
	assertSequence(t, 2, Sequence(3, 16), 16)
	assertSequence(t, 2, Sequence(3, 17), 17)

	assertSequence(t, 0, Sequence(3, 18), 18)
	assertSequence(t, 1, Sequence(3, 21), 21)
	assertSequence(t, 1, Sequence(3, 22), 22)
	assertSequence(t, 2, Sequence(3, 24), 24)

}

func assertSequence(t *testing.T, expected, have, ind int) {
	if expected != have {
		t.Errorf("Bad Sequence, expected %d but have %d for ind %d", expected, have, ind)
	}
}
