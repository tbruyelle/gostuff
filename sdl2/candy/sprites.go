package main

type Sprite struct {
	_type SpriteType
	frame int
}

type SpriteType int

const (
	NoSprite SpriteType = iota
	RedSprite
	GreenSprite
	BlueSprite
	YellowSprite
	PinkSprite
	OrangeSprite
	RedHStripesSprite
	GreenHStripesSprite
	BlueHStripesSprite
	YellowHStripesSprite
	PinkHStripesSprite
	OrangeHStripesSprite
	RedVStripesSprite
	GreenVStripesSprite
	BlueVStripesSprite
	YellowVStripesSprite
	PinkVStripesSprite
	OrangeVStripesSprite
	RedPackedSprite
	GreenPackedSprite
	BluePackedSprite
	YellowPackedSprite
	PinkPackedSprite
	OrangePackedSprite
	BombSprite
	DyingSprite
)

func NewSprite(_type SpriteType) Sprite {
	return Sprite{_type: _type}
}

func (c *Candy) determineIdleSprite() Sprite {
	return NewSprite(SpriteType(c._type))
}
