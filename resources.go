package pctk

// ResourceLocator is the name of a resource.
type ResourceLocator string

// RootLocator is the locator of the root of the resources.
const RootLocator ResourceLocator = ""

// Append appends the other locator to the locator l.
func (l ResourceLocator) Append(other ResourceLocator) ResourceLocator {
	return l + "/" + other
}

// ResourceLoader is a value that can load game resources.
type ResourceLoader interface {
	// LoadActor loads an actor from the given locator. It returns nil if the actor is not found.
	LoadActor(locator ResourceLocator) *Actor

	// LoadMusic loads a music song from the given locator. It returns nil if the music is not found.
	LoadMusic(locator ResourceLocator) *Music

	// LoadScene loads a scene from the given locator. It returns nil if the scene is not found.
	LoadScene(locator ResourceLocator) *Scene

	// LoadScript loads a script from the given locator. It returns nil if the script is not found.
	LoadScript(locator ResourceLocator) *Script

	// LoadSound loads a sound effect from he given locator. It returns nil if the sound is not found.
	LoadSound(locator ResourceLocator) *Sound

	// LoadSpriteSheet loads a sprite sheet from the given locator. It returns nil if the sprite
	// sheet is not found.
	LoadSpriteSheet(locator ResourceLocator) *SpriteSheet
}

// ResourceBundle is a bundle of resources that are loaded in memory. This can be used for
// testing purposes mainly.
type ResourceBundle struct {
	actors       map[ResourceLocator]*Actor
	music        map[ResourceLocator]*Music
	scenes       map[ResourceLocator]*Scene
	scripts      map[ResourceLocator]*Script
	sounds       map[ResourceLocator]*Sound
	spriteSheets map[ResourceLocator]*SpriteSheet
}

// NewResourceBundle creates a new resource bundle that can be used as resource loader.
func NewResourceBundle() *ResourceBundle {
	return &ResourceBundle{
		actors:       make(map[ResourceLocator]*Actor),
		music:        make(map[ResourceLocator]*Music),
		scenes:       make(map[ResourceLocator]*Scene),
		scripts:      make(map[ResourceLocator]*Script),
		sounds:       make(map[ResourceLocator]*Sound),
		spriteSheets: make(map[ResourceLocator]*SpriteSheet),
	}
}

// PutActor adds an actor to the catalog.
func (c *ResourceBundle) PutActor(loc ResourceLocator, a *Actor) {
	c.actors[loc] = a
}

// PutMusic adds a music song to the catalog.
func (c *ResourceBundle) PutMusic(loc ResourceLocator, m *Music) {
	c.music[loc] = m
}

// PutScene adds a scene to the catalog.
func (c *ResourceBundle) PutScene(loc ResourceLocator, sc *Scene) {
	c.scenes[loc] = sc
}

// PutScript adds a script to the catalog.
func (c *ResourceBundle) PutScript(loc ResourceLocator, s *Script) {
	c.scripts[loc] = s
}

// PutSound adds a sound to the catalog.
func (c *ResourceBundle) PutSound(loc ResourceLocator, s *Sound) {
	c.sounds[loc] = s
}

// PutSpriteSheet adds a sprite sheet to the catalog.
func (c *ResourceBundle) PutSpriteSheet(loc ResourceLocator, ss *SpriteSheet) {
	c.spriteSheets[loc] = ss
}

// LoadActor loads an actor from the given locator. It returns nil if the actor is not found.
func (c *ResourceBundle) LoadActor(locator ResourceLocator) *Actor {
	return c.actors[locator]
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

// LoadSpriteSheet loads a sprite sheet from the given locator. It returns nil if the sprite sheet
// is not found.
func (c *ResourceBundle) LoadSpriteSheet(locator ResourceLocator) *SpriteSheet {
	return c.spriteSheets[locator]
}

// LoadSound loads a sound effect from he given locator. It returns nil if the sound is not found.
func (c *ResourceBundle) LoadSound(locator ResourceLocator) *Sound {
	return c.sounds[locator]
}
