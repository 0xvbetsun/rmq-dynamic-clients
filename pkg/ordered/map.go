// Package ordered contains of ordered map structure and methods
package ordered

import "container/list"

// Map stores data in ordered way, where key is equal to value
type Map struct {
	kv map[interface{}]*list.Element
	ll *list.List
}

// NewMap creates new instance of Map
func NewMap() *Map {
	return &Map{
		kv: make(map[interface{}]*list.Element),
		ll: list.New(),
	}
}

// Get returns the value for a key.
func (m *Map) Get(key interface{}) interface{} {
	if el, ok := m.kv[key]; ok {
		return el.Value
	}

	return nil
}

// Add pushes new element to a map. If element has already exists it will return ald value
func (om *Map) Add(key interface{}) interface{} {
	el, found := om.kv[key]
	if !found {
		element := om.ll.PushBack(key)
		om.kv[key] = element
		return element.Value
	}
	return el.Value
}

// Len returns the number of elements in the map.
func (m *Map) Len() int {
	return len(m.kv)
}

// Keys returns all of the keys in the order they were inserted.
func (m *Map) Keys() (keys []interface{}) {
	keys = make([]interface{}, m.Len())

	element := m.ll.Front()
	for i := 0; element != nil; i++ {
		keys[i] = element.Value
		element = element.Next()
	}

	return keys
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (m *Map) Delete(key interface{}) bool {
	element, ok := m.kv[key]
	if ok {
		m.ll.Remove(element)
		delete(m.kv, key)
	}

	return ok
}
