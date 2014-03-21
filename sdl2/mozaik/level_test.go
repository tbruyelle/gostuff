package main

import (
"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsPlain(t *testing.T) {
	lvl:=LoadLevel(1)

assert.True(t, lvl.IsPlain(0))
assert.False(t, lvl.IsPlain(1))
assert.True(t, lvl.IsPlain(2))
}
