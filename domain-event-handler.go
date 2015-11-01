package domain_events

import (
	"fmt"
)

type DomainHandler interface {
	Handle(messageType string, data []byte)
	ListMessageTypes() []string
	AddEventHandler(key string, handler DomainEventHandler) (err error)
}

type Handler struct {
	handlers map[string]DomainEventHandler
	keys     []string
}

func NewHandler() DomainHandler {
	handlers := make(map[string]DomainEventHandler)

	keys := make([]string, 0, len(handlers))
	for key := range handlers {
		keys = append(keys, key)
	}

	return &Handler{handlers: handlers, keys: keys}
}

func (h *Handler) AddEventHandler(key string, handler DomainEventHandler) (err error) {
	h.handlers[key] = handler

	found := false

	for _, currentKey := range h.keys {
		if currentKey == key {
			found = true
		}
	}

	if !found {
		h.keys = append(h.keys, key)
	}

	return nil
}

func (h *Handler) ListMessageTypes() (messageTypes []string) {
	return h.keys
}

func (h *Handler) Handle(messageType string, data []byte) {
	if msgHandler := h.handlers[messageType]; msgHandler != nil {
		fmt.Printf("[HANDLER] | %s - %s\n", messageType, data)
		msgHandler(data)
	}
}

type DomainEventHandler func(eventJson []byte)
