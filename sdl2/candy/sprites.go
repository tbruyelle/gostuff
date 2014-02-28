package main

type Sprite struct {
	_type SpriteType
	frame int
}

type SpriteType int

const (
	// CandySprite represents the sprite according
	// to the candy type.
	CandySprite SpriteType = iota
	DyingSprite
)

func NewSprite(_type SpriteType) Sprite {
	return Sprite{_type: _type}
}
