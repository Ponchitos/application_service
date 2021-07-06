package pkg

import "context"

type Publisher interface {
	Publish(topic string, messages ...interface{}) error

	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan interface{}, error)

	Close() error
}
