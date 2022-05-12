// Package storage represents manipulation with data in app
package storage

import (
	"log"

	"github.com/vbetsun/rmq-dynamic-clients/pkg/ordered"
)

// Items stores ordered map and uses for rpc handling
type Items struct {
	*ordered.Map
}

// NewItems creates new instans of Items
func NewItems() *Items {
	return &Items{
		ordered.NewMap(),
	}
}

// AddItem adds received item to the ordered map
func (i *Items) AddItem(val string, item *interface{}) error {
	*item = i.Add(val)
	log.Printf("add item: %s", *item)
	return nil
}

// GetItem returns requested item
func (i *Items) GetItem(val string, item *interface{}) error {
	*item = i.Get(val)
	log.Printf("get item %s, %s", val, *item)
	return nil
}

// GetAllItems returns items in order they were added
func (i *Items) GetAllItems(val string, items *[]interface{}) error {
	*items = i.Keys()
	log.Printf("get all items: %q", *items)
	return nil
}

// RemoveItem deletes item from ordered map and keep order of items
func (i *Items) RemoveItem(val string, deleted *bool) error {
	*deleted = i.Delete(val)
	log.Printf("item %s was removed %t", val, *deleted)
	return nil
}
