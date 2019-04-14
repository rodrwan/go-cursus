package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/rodrwan/go-cursus/emitter"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var topic = flag.String("topic", "users", "topic to subscribe")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	emit, err := emitter.New(*topic)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			message := t.String()
			emit.Emit("create", message)
		case <-interrupt:
			emit.Disconnect()
			return
		}
	}
}
