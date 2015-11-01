package domain_events

import (
	"fmt"
	"log"
)

type Publisher interface {
	PublishMessage(message DomainEvent) (err error)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
