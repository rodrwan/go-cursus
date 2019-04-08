package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Finciero/cursus"
	"github.com/Finciero/cursus/service"
	"github.com/google/uuid"
)

func main() {
	svc := service.Cursus{}
	svc.Init()
	p := cursus.NewPublisher(service.HelloWorld)
	ud := cursus.NewPublisher(service.UpdateDollar)

	svc.AddPublisher(service.HelloWorld, p)
	svc.AddPublisher(service.UpdateDollar, ud)
	svc.Run()

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		defer r.Body.Close()

		var req cursus.SubscriptionRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		svc.AddSubscriber(req.Topic, &service.Subscriber{
			ID: uuid.New().String(),
		})

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		defer r.Body.Close()

		var req cursus.PublishRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		m := &cursus.Message{
			Data:      req.Message,
			Timestamp: time.Now(),
		}

		svc.Emit(req.Topic, m)
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
