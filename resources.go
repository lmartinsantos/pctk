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
	// LoadScene loads a scene from the given locator. It returns nil if the scene is not found.
	LoadScene(locator ResourceLocator) *Scene

	// LoadSpriteSheet loads a sprite sheet from the given locator. It returns nil if the sprite
	// sheet is not found.
	LoadSpriteSheet(locator ResourceLocator) *SpriteSheet

	// LoadActor loads an actor from the given locator. It returns nil if the actor is not found.
	LoadActor(locator ResourceLocator) *Actor

	// LoadMusic loads a music song from the given locator. It returns nil if the music is not found.
	LoadMusic(locator ResourceLocator) *Music
}

// ResourceBundle is a bundle of resources that are loaded in memory. This can be used for
// testing purposes mainly.
type ResourceBundle struct {
	scenes       map[ResourceLocator]*Scene
	spriteSheets map[ResourceLocator]*SpriteSheet
	actors       map[ResourceLocator]*Actor
	songs        map[ResourceLocator]*Music
}

// NewResourceBundle creates a new resource bundle that can be used as resource loader.
func NewResourceBundle() *ResourceBundle {
	return &ResourceBundle{
		scenes:       make(map[ResourceLocator]*Scene),
		spriteSheets: make(map[ResourceLocator]*SpriteSheet),
		actors:       make(map[ResourceLocator]*Actor),
		songs:        make(map[ResourceLocator]*Music),
	}
}

// PutScene adds a scene to the catalog.
func (c *ResourceBundle) PutScene(loc ResourceLocator, sc *Scene) {
	c.scenes[loc] = sc
}

// PutSpriteSheet adds a sprite sheet to the catalog.
func (c *ResourceBundle) PutSpriteSheet(loc ResourceLocator, ss *SpriteSheet) {
	c.spriteSheets[loc] = ss
}

// PutActor adds an actor to the catalog.
func (c *ResourceBundle) PutActor(loc ResourceLocator, a *Actor) {
	c.actors[loc] = a
}

// PutMusic adds a music song to the catalog.
func (c *ResourceBundle) PutMusic(loc ResourceLocator, m *Music) {
	c.songs[loc] = m
}

// LoadScene loads a scene from the given locator. It returns nil if the scene is not found.
func (c *ResourceBundle) LoadScene(locator ResourceLocator) *Scene {
	return c.scenes[locator]
}

// LoadSpriteSheet loads a sprite sheet from the given locator. It returns nil if the sprite sheet
// is not found.
func (c *ResourceBundle) LoadSpriteSheet(locator ResourceLocator) *SpriteSheet {
	return c.spriteSheets[locator]
}

// LoadActor loads an actor from the given locator. It returns nil if the actor is not found.
func (c *ResourceBundle) LoadActor(locator ResourceLocator) *Actor {
	return c.actors[locator]
}

// LoadMusic loads a music song from the given locator. It returns nil if the music is not found.
func (c *ResourceBundle) LoadMusic(locator ResourceLocator) *Music {
	return c.songs[locator]
}
