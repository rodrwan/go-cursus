package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Subscriber 2")

	data := []byte(`
{
	"topic": "update-dollar"
}
	`)

	req, err := http.NewRequest("POST", "http://localhost:8080/subscribe", bytes.NewBuffer(data))
	check(err)

	hc := http.Client{}

	resp, err := hc.Do(req)
	check(err)

	if resp.StatusCode != http.StatusOK {
		check(errors.New("Bad request"))
	}

	// go svc.Listen()
	forever := make(chan bool)

	<-forever
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
