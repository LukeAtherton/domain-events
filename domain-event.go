package domain_events

import (
	"github.com/satori/go.uuid"

	"time"
)

type DomainEvent struct {
	*MessageHeader `json:"header" xml:"header"`
	Body           interface{} `json:"body" xml:"body"`
}

// type Event struct {
// 	*MessageHeader `json:"header" xml:"header"`
// 	Body           interface{} `json:"body" xml:"body"`
// }

type MessageHeader struct {
	CorrelationId uuid.UUID    `json:"c_id" xml:"c_id"`
	MessageType   string       `json:"message_type" xml:"message_type"`
	Source        *EventSource `json:"source" xml:"source"`
	SentAt        time.Time    `json:"sent_at" xml:"sent_at"`
}

func (h *MessageHeader) GetHeader() (header *MessageHeader) {
	return h
}

func (h *MessageHeader) GetMessageType() (messageType string) {
	return h.MessageType
}

type EventSource struct {
	Service  string `json:"service" xml:"service"`
	Trigger  string `json:"trigger" xml:"trigger"`
	SenderId string `json:"sender_id" xml:"sender_id"`
}

func BuildHeader(messageType string, source *EventSource) (header *MessageHeader) {
	return &MessageHeader{CorrelationId: uuid.NewV4(), MessageType: messageType, Source: source, SentAt: time.Now().UTC()}
}

func BuildEvent(messageType string, source *EventSource, data interface{}) (event DomainEvent) {
	return DomainEvent{MessageHeader: BuildHeader(messageType, source), Body: data}
}
