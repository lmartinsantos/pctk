package pctk

import "io"

// ResourceLocator is the name of a resource.
type ResourceLocator string

// RootLocator is the locator of the root of the resources.
const RootLocator ResourceLocator = ""

func (l ResourceLocator) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, string(l))
}

// Append appends the other locator to the locator l.
func (l ResourceLocator) Append(other ResourceLocator) ResourceLocator {
	return l + "/" + other
}

// ResourceLoader is a value that can load game resources.
type ResourceLoader interface {
	// LoadCostume loads a costume from the given locator. It returns nil if the costume is not
	// found.
	LoadCostume(locator ResourceLocator) *Costume

	// LoadMusic loads a music song from the given locator. It returns nil if the music is not
	// found.
	LoadMusic(locator ResourceLocator) *Music

	// LoadScene loads a scene from the given locator. It returns nil if the scene is not found.
	LoadScene(locator ResourceLocator) *Scene

	// LoadScript loads a script from the given locator. It returns nil if the script is not found.
	LoadScript(locator ResourceLocator) *Script

	// LoadSound loads a sound effect from he given locator. It returns nil if the sound is not
	// found.
	LoadSound(locator ResourceLocator) *Sound
}

// ResourceBundle is a bundle of resources that are loaded in memory. This can be used for
// testing purposes mainly.
type ResourceBundle struct {
	costumes map[ResourceLocator]*Costume
	music    map[ResourceLocator]*Music
	scenes   map[ResourceLocator]*Scene
	scripts  map[ResourceLocator]*Script
	sounds   map[ResourceLocator]*Sound
}

// NewResourceBundle creates a new resource bundle that can be used as resource loader.
func NewResourceBundle() *ResourceBundle {
	return &ResourceBundle{
		costumes: make(map[ResourceLocator]*Costume),
		music:    make(map[ResourceLocator]*Music),
		scenes:   make(map[ResourceLocator]*Scene),
		scripts:  make(map[ResourceLocator]*Script),
		sounds:   make(map[ResourceLocator]*Sound),
	}
}

// PutCostume adds a costume to the bundle.
func (c *ResourceBundle) PutCostume(loc ResourceLocator, cos *Costume) {
	c.costumes[loc] = cos
}

// PutMusic adds a music song to the bundle.
func (c *ResourceBundle) PutMusic(loc ResourceLocator, m *Music) {
	c.music[loc] = m
}

// PutScene adds a scene to the bundle.
func (c *ResourceBundle) PutScene(loc ResourceLocator, sc *Scene) {
	c.scenes[loc] = sc
}

// PutScript adds a script to the bundle.
func (c *ResourceBundle) PutScript(loc ResourceLocator, s *Script) {
	c.scripts[loc] = s
}

// PutSound adds a sound to the bundle.
func (c *ResourceBundle) PutSound(loc ResourceLocator, s *Sound) {
	c.sounds[loc] = s
}

// LoadCostume loads a costume from the given locator. It returns nil if the costume is not found.
func (c *ResourceBundle) LoadCostume(locator ResourceLocator) *Costume {
	return c.costumes[locator]
}

// LoadMusic loads a music song from the given locator. It returns nil if the music is not found.
func (c *ResourceBundle) LoadMusic(locator ResourceLocator) *Music {
	return c.music[locator]
}

// LoadScene loads a scene from the given locator. It returns nil if the scene is not found.
func (c *ResourceBundle) LoadScene(locator ResourceLocator) *Scene {
	return c.scenes[locator]
}

// LoadScript loads a script from the given locator. It returns nil if the script is not found.
func (c *ResourceBundle) LoadScript(locator ResourceLocator) *Script {
	return c.scripts[locator]
}

// LoadSound loads a sound effect from he given locator. It returns nil if the sound is not found.
func (c *ResourceBundle) LoadSound(locator ResourceLocator) *Sound {
	return c.sounds[locator]
}
