package main

import (
	"fmt"

	brk "github.com/nanernunes/federation/pkg/brokers"
	fed "github.com/nanernunes/federation/pkg/federation"
)

const Debug = true

func main() {

	brokers := fed.GetBrokers()

	for _, mapping := range fed.GetMappings() {
		go func(m fed.Mapping, brokers map[string]brk.Broker) {
			if Debug {
				fmt.Printf(
					"[ Mapping %s ]: %s:%s -> %s:%s\n",
					m.Name,
					m.Source.Broker,
					m.Source.Topic,
					m.Target.Broker,
					m.Target.Topic,
				)
			}

			messages, _ := brokers[m.Source.Broker].Subscribe(m.Source.Topic)
			for msg := range messages {
				if Debug {
					fmt.Println(string(msg.Body), msg.Headers)
				}
				brokers[m.Source.Broker].Ack(&msg)
			}
		}(mapping, brokers)
	}

	<-make(chan bool)
}
