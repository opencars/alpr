package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/opencars/alpr/pkg/logger"
	"github.com/opencars/alpr/pkg/queue"
)

const (
	topic = "events.alpr.recognized"
)

type Queue struct {
	conn *nats.Conn
}

func New(url string, enabled bool) (*Queue, error) {
	if !enabled {
		return &Queue{}, nil
	}

	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Queue{
		conn: conn,
	}, nil
}

func (p *Queue) Publish(event *queue.Event) error {
	if p.conn == nil {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.conn.Publish(topic, data)
}

func (p *Queue) Subscribe(ctx context.Context) (<-chan *queue.Event, error) {
	messages := make(chan *nats.Msg, 64)
	_, err := p.conn.ChanSubscribe(topic, messages)
	if err != nil {
		return nil, err
	}

	events := make(chan *queue.Event)
	go msgToEvent(ctx, messages, events)

	return events, nil
}

func msgToEvent(ctx context.Context, in <-chan *nats.Msg, out chan<- *queue.Event) {
	iter := func() (bool, error) {
		select {
		case msg, ok := <-in:
			if !ok {
				return false, nil
			}

			var event queue.Event

			err := json.Unmarshal(msg.Data, &event)
			if err != nil {
				return true, err
			}

			out <- &event
		case <-ctx.Done():
			return false, nil
		}

		return true, nil
	}

	for {
		resume, err := iter()
		if err != nil {
			logger.Errorf("nats: iteration: %v", err)
		}

		if !resume {
			logger.Debugf("nats: stopped")
			return
		}
	}
}
