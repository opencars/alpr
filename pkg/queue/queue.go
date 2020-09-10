package queue

import "context"

type Publisher interface {
	Publish(event *Event) error
}

type Subscriber interface {
	Subscribe(ctx context.Context) (<-chan *Event, error)
}

type Queue interface {
	Publisher
	Subscriber
}
