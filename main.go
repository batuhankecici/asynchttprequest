package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PizzaOrder struct {
	Pizza, Store, Price string
}

func main() {

	var pizza = flag.String("pizza", "", "Pizza to order")
	var store = flag.String("store", "", "Name of the Pizza Store")
	var price = flag.String("price", "", "Price")

	flag.Parse()

	order := PizzaOrder{*pizza, *store, *price}
	body, _ := json.Marshal(order)

	start := time.Now()

	orderChan := make(chan *http.Response)
	paymentChan := make(chan *http.Response)
	storeChan := make(chan *http.Response)

	// OrderService is expected at 8081
	go SendPostAsync("http://google.com", body, orderChan)

	// PaymentService is expected at 8082
	go SendPostAsync("http://facebook.com", body, paymentChan)

	// StoreService is expected at 8083
	go SendPostAsync("http://twitter.com", body, storeChan)

	orderResponse := <-orderChan
	defer orderResponse.Body.Close()
	bytes, _ := ioutil.ReadAll(orderResponse.Body)
	fmt.Println(string(bytes))

	paymentResponse := <-paymentChan
	defer paymentResponse.Body.Close()
	bytes, _ = ioutil.ReadAll(paymentResponse.Body)
	fmt.Println(string(bytes))

	storeResponse := <-storeChan
	defer storeResponse.Body.Close()
	bytes, _ = ioutil.ReadAll(storeResponse.Body)
	fmt.Println(string(bytes))

	end := time.Now()

	fmt.Printf("Order processed after %v seconds\n", end.Sub(start).Seconds())
}

// SendPostAsync send a HTTP Post request to the given url an puts the response
// into the given response channel rc
func SendPostAsync(url string, body []byte, rc chan *http.Response) {
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	rc <- response
}

func SendPostRequest(url string, body []byte) *http.Response {
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	return response
}

