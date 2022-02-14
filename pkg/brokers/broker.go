package brokers

type Broker interface {
	Ack(*Message) error
	Publish(string, *Message, map[string]interface{}) (string, error)
	Subscribe(string) (<-chan Message, error)
}
