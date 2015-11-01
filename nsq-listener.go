package domain_events

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/bitly/go-simplejson"
	"github.com/bitly/nsq/util"
)

type NsqListener struct {
	handler                  DomainHandler
	appName                  string
	topic                    string
	nsqLookupAddressHttp     string
	nsqNodeTcp               string
	resultChannel            chan string
	requireJsonValueParsed   bool
	requireJsonValueIsNumber bool
	requireJsonNumber        float64
	listenStart              time.Time
}

func NewMsgQListener(appName string, topic string, nsqLookupAddressHttp string, nsqNodeTcp string) *NsqListener {

	return &NsqListener{appName: appName, handler: NewHandler(), topic: topic, nsqLookupAddressHttp: nsqLookupAddressHttp, nsqNodeTcp: nsqNodeTcp}
}

func (listener *NsqListener) shouldPassMessage(jsonMsg *simplejson.Json) (bool, bool) {
	pass := true
	backoff := false

	var requireJsonField = "message_type"
	var requireJsonValue = "" //handler.correlationId

	if requireJsonField == "" {
		return pass, backoff
	}

	if requireJsonValue != "" && !listener.requireJsonValueParsed {
		// cache conversion in case needed while filtering json
		var err error
		listener.requireJsonNumber, err = strconv.ParseFloat(requireJsonValue, 64)
		listener.requireJsonValueIsNumber = (err == nil)
		listener.requireJsonValueParsed = true
	}

	header, checkErr := jsonMsg.CheckGet("header")

	if checkErr == false {
		log.Printf("ERROR: Unable to get header")
		pass = false
	}

	jsonVal, ok := header.CheckGet(requireJsonField)
	if !ok {
		pass = false
		log.Printf("ERROR: Incorrect value: %s \n", jsonVal)
		/*if requireJsonValue != "" {
			log.Printf("ERROR: missing field to check required value")
			backoff = true
		}*/
	} else if requireJsonValue != "" {
		// if command-line argument can't convert to float, then it can't match a number
		// if it can, also integers (up to 2^53 or so) can be compared as float64
		if strVal, err := jsonVal.String(); err == nil {
			if strVal != requireJsonValue {
				pass = false
			}
		} else if listener.requireJsonValueIsNumber {
			floatVal, err := jsonVal.Float64()
			if err != nil || listener.requireJsonNumber != floatVal {
				pass = false
			}
		} else {
			// json value wasn't a plain string, and argument wasn't a number
			// give up on comparisons of other types
			pass = false
		}
	}

	return pass, backoff
}

func (listener *NsqListener) HandleMessage(m *nsq.Message) error {
	handleTime := time.Now()

	var err error

	var requireJsonField = "header"

	if requireJsonField != "" {
		var jsonMsg *simplejson.Json
		jsonMsg, err = simplejson.NewJson(m.Body)
		//fmt.Printf("Got: %s\n", jsonMsg)
		if err != nil {
			log.Printf("ERROR: Unable to decode json: %s", m.Body)
			return nil
		}

		//check if should handle message
		if pass, _ := listener.shouldPassMessage(jsonMsg); !pass {
			/*if backoff {
				return errors.New("backoff")
			}*/
			log.Println("ERROR: Didn't Pass")
			return nil
		}

		header, checkErr := jsonMsg.CheckGet("header")

		if checkErr == false {
			log.Printf("ERROR: Unable to get header")
			return nil
		}

		messageType, checkErr := header.CheckGet("message_type")

		if checkErr == false {
			log.Printf("ERROR: Unable to get message_type")
			return nil
		}

		if messageTypeName, err := messageType.String(); err == nil {

			log.Printf("Got Message: %s", messageTypeName)

			listener.handler.Handle(messageTypeName, m.Body)

			fmt.Printf("Listened for: %v\n", time.Now().Sub(listener.listenStart))
			fmt.Printf("Handled in: %v\n", time.Now().Sub(handleTime))
		}
	} else {
		log.Printf("ERROR: %s", err.Error())
	}

	return nil
}

func (listener *NsqListener) Listen() {

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	//resultChan := make(chan string, 1)
	//correlationId := NewUUID().String()
	//listenStart := time.Now()

	maxInFlight := 250

	cfg := nsq.NewConfig()
	cfg.Set("max_in_flight", maxInFlight)
	//cfg.Set("max_backoff_duration", 10)
	cfg.UserAgent = fmt.Sprintf("%s go-nsq/%s", listener.appName, nsq.VERSION)

	var readerOpts = util.StringArray{}

	err := util.ParseOpts(cfg, readerOpts)
	if err != nil {
		log.Fatalf(err.Error())
	}

	setupTime := time.Now()

	channel := "saas-users-app" //fmt.Sprintf("%s", h.appName)

	if !util.IsValidTopicName(listener.topic) {
		log.Fatalf("--topic is invalid")
	}

	if !util.IsValidChannelName(channel) {
		log.Fatalf("--channel is invalid")
	}

	q, err := nsq.NewConsumer(listener.topic, channel, cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	q.SetLogger(log.New(os.Stderr, "", log.LstdFlags), nsq.LogLevelInfo)

	q.AddHandler(listener)

	//for _, addrString := range nsqdTCPAddrs {
	err = q.ConnectToNSQD(listener.nsqNodeTcp)
	if err != nil {
		log.Fatalf(err.Error())
	}
	//}

	//for _, addrString := range lookupdHTTPAddrs {
	log.Printf("lookupd addr %s", listener.nsqLookupAddressHttp)
	err = q.ConnectToNSQLookupd(listener.nsqLookupAddressHttp)
	if err != nil {
		log.Fatalf(err.Error())
	}
	//}

	fmt.Printf("Setup in: %v\n", time.Now().Sub(setupTime))

	for {
		select {
		case <-q.StopChan:
			return
		case <-termChan:
			q.Stop()
		}
	}
}
