package federation

import (
	brk "github.com/nanernunes/federation/pkg/brokers"
)

type Endpoint struct {
	Broker brk.Broker `json:"broker"`
	Topic  string     `json:"topic"`
}
