package domain_events

import (
	"encoding/json"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"
)

type pubsubPublisher struct {
	topic     string
	pubSubCtx context.Context
}

func NewPubsubPublisher(projectId string, topic string) Publisher {
	pubSubCtx, err := cloudContext(projectId)

	if err != nil {
		log.Fatal(err)
	}

	publisher := &pubsubPublisher{
		topic:     topic,
		pubSubCtx: pubSubCtx,
	}
	return publisher
}

func (publisher *pubsubPublisher) PublishMessage(message DomainEvent) (err error) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return
	}
	_, err = pubsub.Publish(publisher.pubSubCtx, publisher.topic, &pubsub.Message{
		Data: jsonData,
	})

	return err
}

func cloudContext(projectID string) (context.Context, error) {
	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, pubsub.ScopePubSub)
	if err != nil {
		return nil, err
	}
	return cloud.WithContext(ctx, projectID, httpClient), nil
}
