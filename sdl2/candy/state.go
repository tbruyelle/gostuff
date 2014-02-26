package main

// State exposes the state methods
type State interface {
	Enter(c *Candy)
	Exit(c *Candy)
	Update(c *Candy)
}

type baseState struct{}

func (s *baseState) Enter(c *Candy) {
	c.sprite = c.determineIdleSprite()
}

func (s *baseState) Exit(c *Candy) {
}

func (s *baseState) Update(c *Candy) {
}

type idleState struct {
	baseState
}

type dyingState struct {
	baseState
	beforeDie int
}

func (s *dyingState) Enter(c *Candy) {
	c.sprite = NewSprite(DyingSprite)
}

func (s *dyingState) Update(c *Candy) {
	if s.beforeDie < 10 {
		s.beforeDie++
		c.sprite.frame++
	} else {
		c.dead = true
	}
}

func NewIdleState() State {
	return &idleState{}
}

func NewDyingState() State {
return &dyingState{}
}
