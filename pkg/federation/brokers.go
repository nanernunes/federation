package federation

import (
	"strings"

	"github.com/nanernunes/federation/pkg/brokers"
	"github.com/nanernunes/federation/pkg/brokers/amqp"
	"github.com/nanernunes/federation/pkg/brokers/aws"
	"github.com/nanernunes/federation/pkg/util/env"
)

func GetBrokers() (items map[string]brokers.Broker) {

	indexes := make(map[string]string)
	items = make(map[string]brokers.Broker)

	for _, proto := range []string{"AMQP", "SNS"} {
		for _, env := range env.LookupEnvsByPrefix(proto) {
			names := strings.Split(env, "_")
			indexes[strings.Join(names[0:2], "_")] = proto
		}
	}

	for broker, proto := range indexes {

		switch proto {
		case "AMQP":
			var config amqp.AMQPConfig
			env.Fetch(broker, &config)
			items[strings.Split(broker, "_")[1]] = amqp.NewAMQP(broker, &config)

		case "SNS":
			var config aws.AWSConfig
			env.Fetch(broker, &config)
			items[strings.Split(broker, "_")[1]] = aws.NewSNS(broker, aws.NewAWS(&config))

		default:
			continue
		}

	}

	return
}
