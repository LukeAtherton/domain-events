package domain_events

import (
	. "github.com/lukeatherton/identity"

	"time"
)

type DomainEvent interface {
	GetMessageType() (messageType string)
	GetHeader() (header *MessageHeader)
}

type MessageHeader struct {
	CorrelationId ID           `json:"c_id" xml:"c_id"`
	MessageType   string       `json:"message_type" xml:"message_type"`
	Sender        *EventSource `json:"sender" xml:"sender"`
	TimeStamp     time.Time    `json:"timestamp" xml:"timestamp"`
}

func (h *MessageHeader) GetHeader() (header *MessageHeader) {
	return h
}

func (h *MessageHeader) GetMessageType() (messageType string) {
	return h.MessageType
}

type EventSource struct {
	Service string `json:"service" xml:"service"`
	UserId  ID     `json:"user_id" xml:"user_id"`
}

func BuildHeader(messageType string, sender *EventSource) (header *MessageHeader) {
	return &MessageHeader{CorrelationId: NewUUID(), MessageType: messageType, Sender: sender, TimeStamp: time.Now().UTC()}
}
