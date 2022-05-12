package rpc

type Service interface {
	AddItem(val string, item *interface{}) error
	GetItem(val string, item *interface{}) error
	GetAllItems(val string, items *[]interface{}) error
	RemoveItem(val string, deleted *bool) error
}

// Handler represents rest modules of API
type Handler struct {
	service Service
}

// NewHandler returns instance of rpc handler
func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// AddItem sends request over rpc for adding item to the store
func (h *Handler) AddItem(val string, item *interface{}) error {
	return h.service.AddItem(val, item)
}

// GetItem sends request over rpc for retrieving item from store
func (h *Handler) GetItem(val string, item *interface{}) error {
	return h.service.GetItem(val, item)
}

// GetAllItems sends request over rpc for retrieving all items in the order they were added
func (h *Handler) GetAllItems(val string, items *[]interface{}) error {
	return h.service.GetAllItems(val, items)
}

// RemoveItem sends request over rpc for deleting requested item from store
func (h *Handler) RemoveItem(val string, deleted *bool) error {
	return h.service.RemoveItem(val, deleted)
}
