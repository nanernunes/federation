package brokers

type Message struct {
	Original interface{}

	Body    string
	Headers map[string]interface{}
}
