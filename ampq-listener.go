package domain_events

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

type AmpqListener struct {
	handler    DomainHandler
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
}

func NewAmpqListener(topic string, queue string, address string, username string, password string, handler DomainHandler) *AmpqListener {

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", username, password, address))
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		topic,   // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for _, s := range handler.ListMessageTypes() {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, topic, s)
		err = ch.QueueBind(
			q.Name, // queue name
			s,      // routing key
			topic,  // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	return &AmpqListener{handler: handler, connection: conn, channel: ch, queue: &q}
}

func (listener *AmpqListener) Listen() {

	defer listener.connection.Close()
	defer listener.channel.Close()

	msgs, err := listener.channel.Consume(
		listener.queue.Name, // queue
		"",                  // consumer
		true,                // auto ack
		false,               // exclusive
		false,               // no local
		false,               // no wait
		nil,                 // args
	)
	failOnError(err, "Failed to register a consumer")

	// forever := make(chan bool)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range msgs {
			// log.Printf(" [x] %s", d.Body)

			listener.handler.Handle(d.RoutingKey, d.Body)
		}
	}()

	log.Printf(" [*] Listening...")
	<-termChan
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
