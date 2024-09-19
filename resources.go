package pctk

import "io"

// ResourceRef is a reference to a resource. This is typically used from resources to refer to other
// resources.
type ResourceRef string

func (r ResourceRef) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, string(r))
}

// Append appends other to reference r.
func (r ResourceRef) Append(other ResourceRef) ResourceRef {
	return r + "/" + other
}

// ResourceLoader is a value that can load game resources.
type ResourceLoader interface {
	// LoadCostume loads a costume from the given ref. It returns nil if the costume is not
	// found.
	LoadCostume(ref ResourceRef) *Costume

	// LoadMusic loads a music song from the given ref. It returns nil if the music is not
	// found.
	LoadMusic(ref ResourceRef) *Music

	// LoadScene loads a scene from the given ref. It returns nil if the scene is not found.
	LoadScene(ref ResourceRef) *Scene

	// LoadScript loads a script from the given ref. It returns nil if the script is not found.
	LoadScript(ref ResourceRef) *Script

	// LoadSound loads a sound effect from he given ref. It returns nil if the sound is not
	// found.
	LoadSound(ref ResourceRef) *Sound
}

// ResourceBundle is a bundle of resources that are loaded in memory. This can be used for
// testing purposes mainly.
type ResourceBundle struct {
	costumes map[ResourceRef]*Costume
	music    map[ResourceRef]*Music
	scenes   map[ResourceRef]*Scene
	scripts  map[ResourceRef]*Script
	sounds   map[ResourceRef]*Sound
}

// NewResourceBundle creates a new resource bundle that can be used as resource loader.
func NewResourceBundle() *ResourceBundle {
	return &ResourceBundle{
		costumes: make(map[ResourceRef]*Costume),
		music:    make(map[ResourceRef]*Music),
		scenes:   make(map[ResourceRef]*Scene),
		scripts:  make(map[ResourceRef]*Script),
		sounds:   make(map[ResourceRef]*Sound),
	}
}

// PutCostume adds a costume to the bundle.
func (c *ResourceBundle) PutCostume(ref ResourceRef, cos *Costume) {
	c.costumes[ref] = cos
}

// PutMusic adds a music song to the bundle.
func (c *ResourceBundle) PutMusic(ref ResourceRef, m *Music) {
	c.music[ref] = m
}

// PutScene adds a scene to the bundle.
func (c *ResourceBundle) PutScene(ref ResourceRef, sc *Scene) {
	c.scenes[ref] = sc
}

// PutScript adds a script to the bundle.
func (c *ResourceBundle) PutScript(ref ResourceRef, s *Script) {
	c.scripts[ref] = s
}

// PutSound adds a sound to the bundle.
func (c *ResourceBundle) PutSound(ref ResourceRef, s *Sound) {
	c.sounds[ref] = s
}

// LoadCostume loads a costume from the given ref. It returns nil if the costume is not found.
func (c *ResourceBundle) LoadCostume(ref ResourceRef) *Costume {
	return c.costumes[ref]
}

// LoadMusic loads a music song from the given ref. It returns nil if the music is not found.
func (c *ResourceBundle) LoadMusic(ref ResourceRef) *Music {
	return c.music[ref]
}

// LoadScene loads a scene from the given ref. It returns nil if the scene is not found.
func (c *ResourceBundle) LoadScene(ref ResourceRef) *Scene {
	return c.scenes[ref]
}

// LoadScript loads a script from the given ref. It returns nil if the script is not found.
func (c *ResourceBundle) LoadScript(ref ResourceRef) *Script {
	return c.scripts[ref]
}

// LoadSound loads a sound effect from he given ref. It returns nil if the sound is not found.
func (c *ResourceBundle) LoadSound(ref ResourceRef) *Sound {
	return c.sounds[ref]
}
