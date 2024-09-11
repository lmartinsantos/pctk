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
}

// ResourceBundle is a bundle of resources that are loaded in memory. This can be used for
// testing purposes mainly.
type ResourceBundle struct {
	scenes map[ResourceLocator]*Scene
}

// NewResourceBundle creates a new resource bundle that can be used as resource loader.
func NewResourceBundle() *ResourceBundle {
	return &ResourceBundle{
		scenes: make(map[ResourceLocator]*Scene),
	}
}

// PutScene adds a scene to the catalog.
func (c *ResourceBundle) PutScene(loc ResourceLocator, sc *Scene) {
	c.scenes[loc] = sc
}

// LoadScene loads a scene from the given locator. It returns nil if the scene is not found.
func (c *ResourceBundle) LoadScene(locator ResourceLocator) *Scene {
	return c.scenes[locator]
}
