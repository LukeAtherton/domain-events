package domain_events

// import (
// 	"crypto/tls"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"strings"
// )

// type nsqHttpPublisher struct {
// 	nsqdHttpAddress string
// }

// func NewPublisher(nsqdHttpAddress string) Publisher {
// 	nsqPublisher := &nsqHttpPublisher{
// 		nsqdHttpAddress: nsqdHttpAddress,
// 	}
// 	return nsqPublisher
// }

// func (publisher *nsqHttpPublisher) PublishMessage(message DomainEvent) (err error) {

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

// 	fmt.Println("Request to User Service: ", urlStr)

// 	resp, err := client.Do(req)

// 	fmt.Println(resp.Status)

// 	if err == nil {
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
