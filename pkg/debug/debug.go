package debug

import (
	"os"
)

func Enabled() bool {
	_, ok := os.LookupEnv("DEBUG")
	return ok
}
