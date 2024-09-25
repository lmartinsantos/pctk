package pctk

// Inventory represents a collection of InventoryItem objects owned by an actor.
type Inventory struct {
	items map[string]*InventoryItem
}

// InventoryItem represents an item in the inventory.
type InventoryItem struct {
	object *Object
	row    int
}

// NewInventory creates and returns a new empty Inventory.
func NewInventory() *Inventory {
	return &Inventory{
		items: make(map[string]*InventoryItem),
	}
}

// AddItem adds an Object to the inventory and assigns it a row based on current items.
func (i *Inventory) AddItem(o *Object) {
	newRow := len(i.items)
	i.items[o.name] = &InventoryItem{object: o, row: newRow}

}

// RemoveItemByName removes an item from the inventory by its name and updates row positions.
func (i *Inventory) RemoveItemByName(name string) {
	item, exists := i.items[name]
	if !exists {
		return
	}

	delete(i.items, name)

	for _, inventoryItem := range i.items {
		if inventoryItem.row > item.row {
			inventoryItem.row--
		}
	}
}

// Bounds is implemented to satisfy the Interactable interface.
func (i *InventoryItem) Bounds() Rectangle {
	x := 2 + 4*ScreenWidth/6
	y := ViewportHeight + (i.row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	return NewRect(x, y, w, h)
}

// Description is implemented to satisfy the Interactable interface.
func (i *InventoryItem) Description() string {
	return i.object.Description()
}
