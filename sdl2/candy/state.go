package main

type StateType int

const (
	IdleStateType StateType = iota
	DyingStateType
)

// State exposes the state methods
type State interface {
	Enter(c *Candy)
	Exit(c *Candy)
	Update(c *Candy)
	Type() StateType
}

type baseState struct{}

func (s *baseState) Enter(c *Candy) {
	c.sprite = NewSprite(CandySprite)
}

func (s *baseState) Exit(c *Candy) {
}

func (s *baseState) Update(c *Candy) {
}

type idleState struct {
	baseState
}

func (s *idleState) Type() StateType{
	return IdleStateType
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

func (s *dyingState) Type() StateType {
	return DyingStateType
}

func NewIdleState() State {
	return &idleState{}
}

func NewDyingState() State {
	return &dyingState{}
}
