package amqp

import (
	"fmt"
	"log"

	brk "github.com/nanernunes/federation/pkg/brokers"

	"github.com/streadway/amqp"
)

type AMQPConfig struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Vhost string
	TLS   bool
}

type AMQP struct {
	Config  AMQPConfig
	Client  *amqp.Connection
	Channel *amqp.Channel
}

func NewAMQP(config *AMQPConfig) *AMQP {
	proto := "amqp"
	if config.TLS {
		proto += "s"
	}

	vhost := config.Vhost
	if vhost == "/" {
		vhost = "%2F"
	}

	conn, err := amqp.Dial(fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		proto, config.User, config.Pass, config.Host, config.Port, vhost,
	))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &AMQP{Client: conn}
}

func (a *AMQP) GetChannel() *amqp.Channel {
	if a.Channel != nil {
		return a.Channel
	}

	ch, err := a.Client.Channel()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return ch
}

func (a *AMQP) Ack(message *brk.Message) error {
	return message.Original.(amqp.Delivery).Ack(false)
}

func (a *AMQP) Subscribe(source string) (<-chan brk.Message, error) {
	queue, err := a.GetChannel().QueueDeclare(
		fmt.Sprintf("%s-%s", source, "federation"),
		true,  // durable
		false, // autoDelete
		false, // exclusive
		true,  // noWait
		amqp.Table{},
	)
	if err != nil {
		fmt.Println(err)
	}

	err = a.GetChannel().QueueBind(
		queue.Name,
		"",     // routingKey
		source, // exchange
		true,   // noWait
		amqp.Table{},
	)
	if err != nil {
		fmt.Println(err)
	}

	deliveryChan, error := a.GetChannel().Consume(
		queue.Name,
		"federation", // consumer
		false,        // autoAck
		false,        // exclusive
		false,        // noLocal
		true,         // noWait
		amqp.Table{},
	)

	messages := make(chan brk.Message)

	go func() {
		for delivery := range deliveryChan {
			headers := make(map[string]interface{})

			for key, value := range delivery.Headers {
				headers[key] = value
			}

			messages <- brk.Message{
				Original: delivery,
				Body:     string(delivery.Body),
				Headers:  delivery.Headers,
			}
		}
	}()

	return messages, error
}

func (a *AMQP) Publish(
	target string, message *brk.Message, options map[string]interface{},
) (string, error) {

	var routingKey string
	if v := options["routingKey"]; v != "" {
		routingKey = v.(string)
	}

	attributes := amqp.Table(message.Headers)

	err := a.GetChannel().Publish(target, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message.Body),
		Headers:     attributes,
	})

	return "", err
}
