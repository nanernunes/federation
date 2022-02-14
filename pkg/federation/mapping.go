package federation

import (
	"os"
	"strings"

	"github.com/nanernunes/federation/pkg/util/env"
)

type Mapping struct {
	Name   string   `json:"name"`
	Source Endpoint `json:"source"`
	Target Endpoint `json:"target"`
}

func GetMappings() (mappings []Mapping) {
	// Ex.: FEDERATION_HELLO_WORLD=MYTARGET_helloworld,MYORIGIN_helloworld

	for _, env := range env.LookupEnvsByPrefix("FEDERATION") {
		name := strings.SplitN(env, "_", 2)[1]        // HELLO_WORLD
		mapping := strings.Split(os.Getenv(env), ",") // MYTARGET_helloworld,MYORIGIN_helloworld

		mappings = append(mappings, Mapping{
			Name: name,
			Source: Endpoint{
				Broker: strings.SplitN(mapping[1], "_", 2)[0], // MYORIGIN
				Topic:  strings.SplitN(mapping[1], "_", 2)[1], // helloworld
			},
			Target: Endpoint{
				Broker: strings.SplitN(mapping[0], "_", 2)[0], // MYTARGET
				Topic:  strings.SplitN(mapping[0], "_", 2)[1], // helloworld
			},
		})
	}
	return
}
