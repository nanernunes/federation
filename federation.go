package main

import (
	fed "github.com/nanernunes/federation/pkg/federation"
)

func main() {

	for _, mapping := range fed.GetMappings() {
		go func(m fed.Mapping) { m.StartForwarding() }(mapping)
	}

	<-make(chan bool)
}
