package pctk

import "io"

// CostumeAction is a value that represents an action for a costume. For predefined actions idle,
// speak, and walk, use the CustomIdle, CustomSpeak, and CustomWalk functions respectively to refer
// to them. For custom actions, use any custom byte value above 0x80.
type CostumeAction byte

// CostumeIdle returns a costume action for the idle action in the given direction.
func CostumeIdle(dir Direction) CostumeAction {
	return CostumeAction((0 << 2) | (dir & 0x03))
}

// CostumeSpeak returns a costume action for the speak action in the given direction.
func CostumeSpeak(dir Direction) CostumeAction {
	return CostumeAction((1 << 2) | (dir & 0x03))
}

// CostumeWalk returns a costume action for the walk action in the given direction.
func CostumeWalk(dir Direction) CostumeAction {
	return CostumeAction((2 << 2) | (dir & 0x03))
}

// Costume is a struct that represents a costume for an actor or a room animation.
type Costume struct {
	sprites *SpriteSheet

	anims map[CostumeAction]*Animation
}

// NewCostume creates a new costume.
func NewCostume(sprites *SpriteSheet) *Costume {
	return &Costume{
		sprites: sprites,
		anims:   make(map[CostumeAction]*Animation),
	}
}

// WithAnimationStand sets the stand animation for the given direction.
func (c *Costume) WithAnimation(act CostumeAction, anim *Animation) *Costume {
	c.anims[act] = anim
	return c
}

// BinaryEncode encodes the costume to a binary format. The format is as follows:
// - uint32: the number of animations.
// - for each animation:
//   - byte: the action.
//   - the animation.
func (c *Costume) BinaryEncode(w io.Writer) (n int, err error) {
	n, err = BinaryEncode(w, uint32(len(c.anims)))
	for act, anim := range c.anims {
		nn, err := BinaryEncode(w, byte(act), anim)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (c *Costume) draw(act CostumeAction, pos Position) {
	if anim := c.anims[act]; anim != nil {
		anim.draw(c.sprites, pos)
	}
}
