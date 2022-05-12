package service

type Storage interface {
	AddItem(val string, item *interface{}) error
	GetItem(val string, item *interface{}) error
	GetAllItems(val string, items *[]interface{}) error
	RemoveItem(val string, deleted *bool) error
}

type Service struct {
	storage Storage
}

// NewService creates an instance of service
func NewService(storage Storage) *Service {
	return &Service{storage}
}

// AddItem sends request over rpc for adding item to the store
func (s Service) AddItem(val string, item *interface{}) error {
	return s.storage.AddItem(val, item)
}

// GetItem sends request over rpc for retrieving item from store
func (s Service) GetItem(val string, item *interface{}) error {
	return s.storage.GetItem(val, item)
}

// GetAllItems sends request over rpc for retrieving all items in the order they were added
func (s Service) GetAllItems(val string, items *[]interface{}) error {
	return s.storage.GetAllItems(val, items)
}

// RemoveItem sends request over rpc for deleting requested item from store
func (s Service) RemoveItem(val string, deleted *bool) error {
	return s.storage.RemoveItem(val, deleted)
}
