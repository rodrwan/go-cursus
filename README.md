# Cursus

Internal message service, to inter connect services to an specific topic.

# Run

To run the broker we need to do the follow:

```sh
$ go run cmd/broker/main.go
```

# Concept behind this implementation

- **topic**: A Topic is the name for a Room.
- **action**: An action is an special keyword that emitter and receiver understand. An action is used to specify special action of your system, like, create, update or delete any record. (See example folder for clarifications).

## Room

A room is labeled by a Topic, on this room we store every client that subscribe to an specific topic.


## Emitter

An Emitter is a struct that let us publish message to an specific `action`.


## Receiver

A Receiver is a struct that listen for incomming message on an specific `Topic`.


## Broker

A broker is an http server that handle websocket connection and let us connect Emitter with Receivers.

