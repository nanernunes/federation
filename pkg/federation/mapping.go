package federation

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nanernunes/federation/pkg/debug"
	"github.com/nanernunes/federation/pkg/util/env"
)

type Mapping struct {
	Name   string   `json:"name"`
	Source Endpoint `json:"source"`
	Target Endpoint `json:"target"`
}

func (m *Mapping) StartForwarding() {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	for {

		if ok := m.Source.Broker.Connect(errChan); !ok {
			time.Sleep(time.Second * 5)
			continue
		}

		if debug.Enabled() {
			log.Printf(
				"[ Mapping %s ]: %s:%s -> %s:%s\n",
				m.Name,
				m.Source.Broker.GetName(),
				m.Source.Topic,
				m.Target.Broker.GetName(),
				m.Target.Topic,
			)
		}

		go func() {
			for msg := range m.Source.Broker.Subscribe(ctx, m.Source.Topic, errChan) {
				if debug.Enabled() {
					log.Printf(string(msg.Body), msg.Headers, "\n")
				}

				m.Source.Broker.Ack(&msg)
			}
		}()

		<-errChan
		cancel()
	}
}

func GetMappings() (mappings []Mapping) {

	brokers := GetBrokers()

	// Ex.: FEDERATION_HELLO_WORLD=MYTARGET_helloworld,MYORIGIN_helloworld
	for _, env := range env.LookupEnvsByPrefix("FEDERATION") {
		name := strings.SplitN(env, "_", 2)[1]        // HELLO_WORLD
		mapping := strings.Split(os.Getenv(env), ",") // MYTARGET_helloworld,MYORIGIN_helloworld

		mappings = append(mappings, Mapping{
			Name: name,
			Source: Endpoint{
				Broker: brokers[strings.SplitN(mapping[1], "_", 2)[0]], // MYORIGIN -> Object
				Topic:  strings.SplitN(mapping[1], "_", 2)[1],          // helloworld
			},
			Target: Endpoint{
				Broker: brokers[strings.SplitN(mapping[0], "_", 2)[0]], // MYTARGET -> Object
				Topic:  strings.SplitN(mapping[0], "_", 2)[1],          // helloworld
			},
		})
	}
	return
}
