package cursus

import (
	"time"
)

// Topic ...
type Topic string

// Message ...
type Message struct {
	Data      string
	Timestamp time.Time
}

// Service ...
type Service interface {
	Run()
	Init()

	AddPublisher(Topic, Publisher)
	AddSubscruber(Topic, Subscriber)
	Emit(Topic, Message)
}

// SubscriptionRequest ...
type SubscriptionRequest struct {
	Topic Topic `json:"topic,omitempty"`
}

// UnsubscriptionRequest ...
type UnsubscriptionRequest struct {
	Topic Topic `json:"topic,omitempty"`
}

// PublishRequest ...
type PublishRequest struct {
	Topic   Topic  `json:"topic,omitempty"`
	Message string `json:"message,omitempty"`
}
