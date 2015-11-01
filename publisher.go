package domain_events

type Publisher interface {
	PublishMessage(message DomainEvent) (err error)
}
