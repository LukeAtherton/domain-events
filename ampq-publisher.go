package domain_events

import (
	// "crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	// "io/ioutil"
	// "net/http"
	// "net/url"
	// "strconv"
	"log"
	// "strings"
)

type ampqPublisher struct {
	uri          string
	exchangeName string
	exchangeType string
	routingKey   string
	reliable     bool
	testMode     bool
}

func NewAmpqPublisher(exchangeAddress string, ampqUsername string, ampqPassword string, topic string) Publisher {
	publisher := &ampqPublisher{
		uri:          fmt.Sprintf("amqp://%s:%s@%s", ampqUsername, ampqPassword, exchangeAddress),
		exchangeName: topic,
		exchangeType: "topic",
		reliable:     true,
		testMode:     false,
	}
	return publisher
}

func (publisher *ampqPublisher) PublishMessage(message DomainEvent) (err error) {

	// This function dials, connects, declares, publishes, and tears down,
	// all in one go. In a real service, you probably want to maintain a
	// long-lived connection as state, and publish against that.

	log.Printf("dialing %q", publisher.uri)
	connection, err := amqp.Dial(publisher.uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	// if err != nil {
	// 	return fmt.Errorf("Dial: %s", err)
	// }
	defer connection.Close()

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	// if err != nil {
	// 	return fmt.Errorf("Channel: %s", err)
	// }

	log.Printf("got Channel, declaring %q Exchange (%q)", publisher.exchangeType, publisher.exchangeName)
	if err := channel.ExchangeDeclare(
		publisher.exchangeName, // name
		publisher.exchangeType, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if publisher.reliable {
		log.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer confirmOne(ack, nack)
	}

	jsonData, _ := json.Marshal(message.Body)
	// post_data := strings.NewReader((string)(jsonData))

	log.Printf("declared Exchange, %s @ %s publishing %s: %dB body (%q)", message.GetHeader().Source.Trigger, message.GetHeader().Source.Service, message.GetHeader().MessageType, len(jsonData), jsonData)
	if err = channel.Publish(
		publisher.exchangeName,          // publish to an exchange
		message.GetHeader().MessageType, // routing to 0 or more queues
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(jsonData),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
			CorrelationId: message.GetHeader().CorrelationId.String(),
			Timestamp:     message.GetHeader().SentAt,
			AppId:         message.GetHeader().Source.Service,
			UserId:        message.GetHeader().Source.SenderId.String(),
			Type:          message.GetHeader().Source.Trigger,
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// func (publisher *rabbitMqPublisher) PublishMessage(message DomainMessage) (err error) {

// 	apiUrl := fmt.Sprintf("http://%s", publisher.nsqdHttpAddress)
// 	resource := "/put"

// 	jsonData, _ := json.Marshal(message)

// 	post_data := strings.NewReader((string)(jsonData))

// 	u, _ := url.ParseRequestURI(apiUrl)
// 	u.Path = resource

// 	//set query params
// 	q := u.Query()
// 	q.Set("topic", "services")
// 	u.RawQuery = q.Encode()

// 	urlStr := fmt.Sprintf("%v", u)

// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}

// 	client := &http.Client{Transport: tr}
// 	req, _ := http.NewRequest("POST", urlStr, post_data) // <-- URL-encoded payload
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Content-Length", strconv.Itoa(post_data.Len()))

// 	fmt.Println("Publish to NSQ: ", urlStr)

// 	resp, err := client.Do(req)

// 	if err == nil {
// 		fmt.Println(resp.Status)

// 		switch resp.StatusCode {
// 		case http.StatusOK:
// 			return nil
// 		case http.StatusCreated:
// 			body, err := ioutil.ReadAll(resp.Body)
// 			if err == nil {
// 				fmt.Println(body)
// 			} else {
// 				panic(err)
// 			}
// 		default:
// 			panic(err)
// 		}
// 	} else {
// 		panic(err)
// 	}

// 	return nil
// }

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(ack, nack chan uint64) {
	log.Printf("waiting for confirmation of one publishing")

	select {
	case tag := <-ack:
		log.Printf("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("failed delivery of delivery tag: %d", tag)
	}
}
