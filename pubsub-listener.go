package domain_events

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	. "google.golang.org/cloud/pubsub"

// 	"golang.org/x/net/context"
// )

// type PubsubListener struct {
// 	handler      DomainHandler
// 	subscription *SubscriptionHandle
// 	pubSubCtx    context.Context
// }

// func NewPubsubPublisher(projectId string, topicName string, subscriptionName string, handler DomainHandler) Listener {
// 	pubSubCtx, err := cloudContext(projectId)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	client, err := NewClient(pubSubCtx, projectId)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	topic := client.Topic(topicName)

// 	doesTopicExist, err := topic.Exists(pubSubCtx)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if !doesTopicExist {
// 		log.Fatal(fmt.Sprintf("ERROR: %s topic does not exist.", topicName))
// 	}

// 	subscription := client.Subscription(subscriptionName)

// 	doesSubExist, err := subscription.Exists(pubSubCtx)

// 	if !doesSubExist {
// 		subscription, err = client.Topic(topicName).Subscribe(pubSubCtx, subscriptionName, time.Duration(0), nil)

// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	listener := &PubsubListener{
// 		handler:      handler,
// 		subscription: subscription,
// 		pubSubCtx:    pubSubCtx,
// 	}

// 	return listener
// }

// // func NewAmpqListener(topic string, queue string, address string, username string, password string, handler DomainHandler) *AmpqListener {

// // 	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// // 	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", username, password, address))
// // 	failOnError(err, "Failed to connect to RabbitMQ")

// // 	ch, err := conn.Channel()
// // 	failOnError(err, "Failed to open a channel")

// // 	err = ch.ExchangeDeclare(
// // 		topic,   // name
// // 		"topic", // type
// // 		true,    // durable
// // 		false,   // auto-deleted
// // 		false,   // internal
// // 		false,   // no-wait
// // 		nil,     // arguments
// // 	)
// // 	failOnError(err, "Failed to declare an exchange")

// // 	q, err := ch.QueueDeclare(
// // 		queue, // name
// // 		false, // durable
// // 		false, // delete when usused
// // 		false, // exclusive
// // 		false, // no-wait
// // 		nil,   // arguments
// // 	)
// // 	failOnError(err, "Failed to declare a queue")

// // 	for _, s := range handler.ListMessageTypes() {
// // 		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, topic, s)
// // 		err = ch.QueueBind(
// // 			q.Name, // queue name
// // 			s,      // routing key
// // 			topic,  // exchange
// // 			false,
// // 			nil)
// // 		failOnError(err, "Failed to bind a queue")
// // 	}

// // 	return &AmpqListener{handler: handler, connection: conn, channel: ch, queue: &q}
// // }

// func (listener *PubsubListener) Listen() {
// 	termChan := make(chan os.Signal, 1)
// 	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		for {
// 			msgs, err := listener.subscription.Pull(listener.pubSubCtx, 0)
// 			if err != nil {
// 				log.Fatalf("could not pull: %v", err)
// 			}

// 			defer msgs.Close()

// 			for m := range msgs {
// 				msg := m

// 				var id int64
// 				if err := json.Unmarshal(msg.Data, &id); err != nil {
// 					log.Printf("could not decode message data: %#v", msg)
// 					go pubsub.Ack(pubSubCtx, subName, msg.AckID)
// 					continue
// 				}

// 				log.Printf("[ID %d] Processing.", id)
// 				go func() {
// 					listener.handler.Handle(d.RoutingKey, d.Body)

// 					pubsub.Ack(pubSubCtx, subName, msg.AckID)
// 					log.Printf("[ID %d] ACK", id)
// 				}()
// 			}
// 		}
// 	}()

// 	log.Printf(" [*] Listening...")
// 	<-termChan
// }

// func (listener *PubsubListener) Listen() {

// 	defer listener.connection.Close()
// 	defer listener.channel.Close()

// 	msgs, err := listener.channel.Consume(
// 		listener.queue.Name, // queue
// 		"",                  // consumer
// 		true,                // auto ack
// 		false,               // exclusive
// 		false,               // no local
// 		false,               // no wait
// 		nil,                 // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	// forever := make(chan bool)
// 	termChan := make(chan os.Signal, 1)
// 	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		for d := range msgs {
// 			// log.Printf(" [x] %s", d.Body)

// 			listener.handler.Handle(d.RoutingKey, d.Body)
// 		}
// 	}()

// 	log.Printf(" [*] Listening...")
// 	<-termChan
// }
