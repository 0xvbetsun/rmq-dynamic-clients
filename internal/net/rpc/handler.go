// Package rpc implements all available requests and handler for communication
package rpc

import (
	"log"
	"net/rpc"

	"github.com/vbetsun/rmq-dynamic-clients/pkg/ordered"
)

// Items stores ordered map and uses for rpc handling
type Items struct {
	om *ordered.Map
}

// NewItems creates new instans of Items
func NewItems() *Items {
	return &Items{
		om: ordered.NewMap(),
	}
}

// AddItem adds received item to the ordered map
func (i *Items) AddItem(val string, item *interface{}) error {
	*item = i.om.Add(val)
	log.Printf("add item: %v", i.om)
	return nil
}

// GetItem returns requested item
func (i *Items) GetItem(val string, item *interface{}) error {
	*item = i.om.Get(val)
	log.Printf("get item: %v", i.om)
	return nil
}

// GetAllItems returns items in order they were added
func (i Items) GetAllItems(val string, items *[]interface{}) error {
	*items = i.om.Keys()
	log.Printf("get all items: %v", i.om)
	return nil
}

// RemoveItem deletes item from ordered map and keep order of items
func (i Items) RemoveItem(val string, deleted *bool) error {
	*deleted = i.om.Delete(val)
	log.Printf("remove item: %v", i.om)
	return nil
}

// AddItem sends request over rpc for adding item to the store
func AddItem(c *rpc.Client, arg string) error {
	return c.Call("Items.AddItem", arg, nil)
}

// GetItem sends request over rpc for retrieving item from store
func GetItem(c *rpc.Client, arg string) error {
	return c.Call("Items.GetItem", arg, nil)
}

// GetAllItems sends request over rpc for retrieving all items in the order they were added
func GetAllItems(c *rpc.Client, arg string) error {
	return c.Call("Items.GetAllItems", arg, nil)
}

// RemoveItem sends request over rpc for deleting requested item from store
func RemoveItem(c *rpc.Client, arg string) error {
	return c.Call("Items.RemoveItem", arg, nil)
}

// BuildRouter creates map for instant finding handler to any command
func BuildRouter() map[string]func(*rpc.Client, string) error {
	return map[string]func(*rpc.Client, string) error{
		"AddItem":     AddItem,
		"RemoveItem":  RemoveItem,
		"GetItem":     GetItem,
		"GetAllItems": GetAllItems,
	}
}
