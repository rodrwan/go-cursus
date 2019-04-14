package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Finciero/cursus/client"
)

var topic = flag.String("topic", "users", "topic to subscribe")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	action, conn, err := client.New(*topic)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case act := <-action:
			log.Printf("message: %s\n", act.Message)
		case <-interrupt:
			fmt.Println("bye bye")
			return
		}
	}
}
