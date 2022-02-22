package amqp

import (
	"context"
	"encoding/json"
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

func (ac *AMQPConfig) GetConnectionString() string {
	proto := "amqp"
	if ac.TLS {
		proto += "s"
	}

	vhost := ac.Vhost
	if vhost == "/" {
		vhost = "%2F"
	}

	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		proto, ac.User, ac.Pass, ac.Host, ac.Port, vhost,
	)
}

type AMQP struct {
	Name    string
	Config  *AMQPConfig
	Client  *amqp.Connection
	Channel *amqp.Channel
	Errors  chan error
}

func NewAMQP(name string, config *AMQPConfig) *AMQP {
	return &AMQP{
		Name:   name,
		Config: config,
		Errors: make(chan error),
	}
}

func (a *AMQP) GetName() string {
	return a.Name
}

func (a *AMQP) Connect(errChan chan error) bool {
	log.Printf("trying to connect %s:%s...\n", a.Config.Host, a.Config.Port)
	conn, err := amqp.Dial(a.Config.GetConnectionString())
	if err != nil {
		return false
	}

	go func() {
		err := <-conn.NotifyClose(make(chan *amqp.Error))
		errChan <- err
	}()

	a.Client = conn
	log.Printf("connection established at %s:%s\n", a.Config.Host, a.Config.Port)
	return true
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

func (a *AMQP) Subscribe(ctx context.Context, source string, chErr chan error) <-chan brk.Message {
	handleError := func(err error) {
		if err != nil {
			chErr <- err
		}
	}

	queue, err := a.GetChannel().QueueDeclare(
		fmt.Sprintf("%s-%s", source, "federation"),
		true,  // durable
		false, // autoDelete
		false, // exclusive
		true,  // noWait
		amqp.Table{},
	)
	handleError(err)

	err = a.GetChannel().QueueBind(
		queue.Name,
		"",     // routingKey
		source, // exchange
		true,   // noWait
		amqp.Table{},
	)
	handleError(err)

	deliveryChan, err := a.GetChannel().Consume(
		queue.Name,
		"federation", // consumer
		false,        // autoAck
		false,        // exclusive
		false,        // noLocal
		true,         // noWait
		amqp.Table{},
	)
	handleError(err)

	messages := make(chan brk.Message)

	go func() {
		for delivery := range deliveryChan {
			headers := make(map[string]interface{})

			for key, value := range delivery.Headers {
				switch value.(type) {
				case amqp.Table:
					if data, err := json.Marshal(value); err != nil {
						headers[key] = string(data)
					}
				default:
					headers[key] = value.(string)
				}
			}

			messages <- brk.Message{
				Original: delivery,
				Body:     string(delivery.Body),
				Headers:  headers,
			}
		}
	}()

	return messages
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

