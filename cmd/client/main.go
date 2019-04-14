package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/Finciero/cursus/receiver"
)

var topic = flag.String("topic", "users", "topic to subscribe")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	receiver, err := receiver.New(*topic)
	if err != nil {
		log.Fatal(err)
	}

	defer receiver.Conn.Close()
	action, err := receiver.Listen()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case act, close := <-action:
			if !close {
				return
			}
			log.Printf("message: %s\n", act.Message)
		case <-interrupt:
			receiver.Disconnect()
			return
		}
	}
}
