package ordered

import (
	"container/list"
)

type Map struct {
	kv map[string]*list.Element
	ll *list.List
}

func NewMap() *Map {
	return &Map{
		kv: make(map[string]*list.Element),
		ll: list.New(),
	}
}

// Get returns the value for a key. If the key does not exist, the second return
// parameter will be false and the value will be nil.
func (ol *Map) Get(key string) string {
	if el, ok := ol.kv[key]; ok {
		return el.Value.(string)
	}

	return ""
}

// Set will set (or replace) a value for a key. If the key was new, then true
// will be returned. The returned value will be false if the value was replaced
// (even if the value was the same).
func (om *Map) Set(key string) string {
	el, found := om.kv[key]
	if !found {
		element := om.ll.PushBack(key)
		om.kv[key] = element
		return element.Value.(string)
	}
	return el.Value.(string)
}

// Len returns the number of elements in the map.
func (ol *Map) Len() int {
	return len(ol.kv)
}

// Keys returns all of the keys in the order they were inserted. If a key was
// replaced it will retain the same position. To ensure most recently set keys
// are always at the end you must always Delete before Set.
func (ol *Map) Keys() (keys []string) {
	keys = make([]string, ol.Len())

	element := ol.ll.Front()
	for i := 0; element != nil; i++ {
		keys[i] = element.Value.(string)
		element = element.Next()
	}

	return keys
}

// Delete will remove a key from the map. It will return true if the key was
// removed (the key did exist).
func (ol *Map) Delete(key string) bool {
	element, ok := ol.kv[key]
	if ok {
		ol.ll.Remove(element)
		delete(ol.kv, key)
	}

	return ok
}
