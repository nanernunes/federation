package brokers

import (
	"context"
)

type Broker interface {
	GetName() string
	Connect(chan error) bool
	Ack(*Message) error
	Publish(string, *Message, map[string]interface{}) (string, error)
	Subscribe(context.Context, string, chan error) <-chan Message
}
