package pctk

import "fmt"

// ResourceLocator is the name of a resource.
type ResourceLocator string

// RootLocator is the locator of the root of the resources.
const RootLocator ResourceLocator = ""

// Append appends the other locator to the locator l.
func (l ResourceLocator) Append(other ResourceLocator) ResourceLocator {
	return l + "/" + other
}

type Resource interface{}

// ResourceCatalog is a catalog of resources that the application can use.
type ResourceCatalog struct {
	index map[ResourceLocator]Resource
}

// NewResourceCatalog creates a new resource catalog.
func NewResourceCatalog() *ResourceCatalog {
	return &ResourceCatalog{
		index: make(map[ResourceLocator]Resource),
	}
}

// Add adds a resource to the catalog.
func (c *ResourceCatalog) Add(loc ResourceLocator, r Resource) {
	c.index[loc] = r
}

// Get gets a resource from the catalog.
func (c *ResourceCatalog) Get(locator ResourceLocator) Resource {
	return c.index[locator]
}

// GetBackground gets the background from the catalog.
func (c *ResourceCatalog) GetBackground(loc ResourceLocator) (*Background, error) {
	bg, ok := c.index[loc].(*Background)
	if !ok {
		return nil, fmt.Errorf("invalid resource type")
	}
	return bg, nil
}
