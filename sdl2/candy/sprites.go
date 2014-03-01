package main

type Sprite struct {
	_type           SpriteType
	nbframes, frame int
}

type SpriteType int

const (
	// CandySprite represents the sprite according
	// to the candy type.
	CandySprite SpriteType = iota
	DyingSprite
)

const (
	StandFrames = 1
	DyingFrames = 10
)

// FramesPerSprites allows to easily grab the number of frames 
// from a SpriteType
var FramesPerSprites = map[SpriteType]int{
	CandySprite: StandFrames,
	DyingSprite: DyingFrames,
}

func NewSprite(_type SpriteType) Sprite {
	return Sprite{_type: _type, nbframes: FramesPerSprites[_type]}
}
