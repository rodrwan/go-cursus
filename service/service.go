package service

import (
	"fmt"
	"time"

	"github.com/Finciero/cursus"
)

// Allowed events
const (
	HelloWorld   cursus.Topic = "hello-world"
	UpdateUser   cursus.Topic = "update-user"
	DeleteUser   cursus.Topic = "delete-user"
	UpdateDollar cursus.Topic = "update-dollar"
)

// Subscriber ...
type Subscriber struct {
	ID string
}

// Do ...
func (s *Subscriber) Do(data string) {
	fmt.Println("subscriber do:", data)
}

type emitMessageToTopic struct {
	Topic   cursus.Topic
	Message *cursus.Message
}

type subscription struct {
	Topic       cursus.Topic
	Subscriptor cursus.Subscriber
}

// Cursus ...
type Cursus struct {
	Subscription       chan *subscription
	Unsubscription     chan *subscription
	EmitMessageToTopic chan *emitMessageToTopic
	Publishers         map[cursus.Topic]*cursus.Publisher
}

// AddPublisher ...
func (c *Cursus) AddPublisher(evt cursus.Topic, p *cursus.Publisher) {
	c.Publishers[evt] = p
}

// AddSubscriber ...
func (c *Cursus) AddSubscriber(evt cursus.Topic, s cursus.Subscriber) {
	c.Subscription <- &subscription{
		Topic:       evt,
		Subscriptor: s,
	}
}

// RemoveSubscriber ...
func (c *Cursus) RemoveSubscriber(evt cursus.Topic) {
	c.Unsubscription <- &subscription{
		Topic: evt,
	}
}

// Init ...
func (c *Cursus) Init() {
	c.Publishers = make(map[cursus.Topic]*cursus.Publisher, 0)
	c.EmitMessageToTopic = make(chan *emitMessageToTopic)
	c.Subscription = make(chan *subscription)
	c.Unsubscription = make(chan *subscription)
}

// Run ...
func (c *Cursus) Run() {
	go func() {
		for {
			select {
			case subs := <-c.Subscription:
				fmt.Println("Subscribe to event:", subs.Topic)
				publisher := c.Publishers[subs.Topic]
				publisher.Subscribe(subs.Topic, subs.Subscriptor)

			case emit := <-c.EmitMessageToTopic:
				fmt.Println("Topic: ", emit.Topic)
				publisher := c.Publishers[emit.Topic]
				for _, subscribes := range publisher.Subscribers {
					for _, subscriber := range subscribes {
						// we should broadcast to everyone on this topic.
						subscriber.Do(emit.Message.Data)
					}
				}
			case <-time.After(time.Millisecond):
			}
		}
	}()
}

// Emit ...
func (c *Cursus) Emit(topic cursus.Topic, m *cursus.Message) {
	c.EmitMessageToTopic <- &emitMessageToTopic{
		Topic:   topic,
		Message: m,
	}
}
