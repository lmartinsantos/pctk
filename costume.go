package pctk

import "io"

// Costume is a struct that represents a costume for an actor or a room animation.
type Costume struct {
	sprites    *SpriteSheet
	animStand  map[Direction]*Animation
	animSpeak  map[Direction]*Animation
	animWalk   map[Direction]*Animation
	animCustom map[string]*Animation
}

// NewCostume creates a new costume.
func NewCostume(sprites *SpriteSheet) *Costume {
	return &Costume{
		sprites:    sprites,
		animStand:  make(map[Direction]*Animation),
		animSpeak:  make(map[Direction]*Animation),
		animWalk:   make(map[Direction]*Animation),
		animCustom: make(map[string]*Animation),
	}
}

// WithAnimationStand sets the stand animation for the given direction.
func (c *Costume) WithAnimationStand(dir Direction, anim *Animation) *Costume {
	c.animStand[dir] = anim
	return c
}

// WithAnimationSpeak sets the speak animation for the given direction.
func (c *Costume) WithAnimationSpeak(dir Direction, anim *Animation) *Costume {
	c.animSpeak[dir] = anim
	return c
}

// WithAnimationWalk sets the walk animation for the given direction.
func (c *Costume) WithAnimationWalk(dir Direction, anim *Animation) *Costume {
	c.animWalk[dir] = anim
	return c
}

// WithAnimationCustom sets a custom animation for the given name.
func (c *Costume) WithAnimationCustom(name string, anim *Animation) *Costume {
	c.animCustom[name] = anim
	return c
}

// BinaryEncode encodes the costume to a binary format.
func (c *Costume) BinaryEncode(w io.Writer) (int, error) {
	panic("not implemented")
}

func (c *Costume) drawStand(pos Position, dir Direction) {
	if anim := c.animStand[dir]; anim != nil {
		anim.draw(c.sprites, pos)
	}
}

func (c *Costume) drawSpeak(pos Position, dir Direction) {
	if anim := c.animSpeak[dir]; anim != nil {
		anim.draw(c.sprites, pos)
	}
}

func (c *Costume) drawWalk(pos Position, dir Direction) {
	if anim := c.animWalk[dir]; anim != nil {
		anim.draw(c.sprites, pos)
	}
}

func (c *Costume) drawCustom(pos Position, name string) {
	if anim := c.animCustom[name]; anim != nil {
		anim.draw(c.sprites, pos)
	}
}
