package cursus

// Publisher ...
type Publisher struct {
	Name        Topic
	Subscribers map[Topic][]Subscriber
}

// Subscribe ...
func (p *Publisher) Subscribe(evt Topic, s Subscriber) {
	p.Subscribers[evt] = append(p.Subscribers[evt], s)
}

// Unsubscribe ...
func (p *Publisher) Unsubscribe(evt Topic) {
	delete(p.Subscribers, evt)
}

// NewPublisher ...
func NewPublisher(topic Topic) *Publisher {
	return &Publisher{
		Name:        topic,
		Subscribers: make(map[Topic][]Subscriber),
	}
}
