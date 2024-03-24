package env

import (
	"encoding/json"
	"expvar"
	"os"

	"github.com/rs/zerolog/log"
)

var env = expvar.NewMap("env")

// Get returns the value of the given environment variable.
func Get(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		value = defaultValue
	}

	log.Debug().Str("name", name).Msg("loaded env var")

	env.Set(name, jsonStringer(value))

	return value
}

type jsonStringer string

func (s jsonStringer) String() string {
	v, _ := json.Marshal(s)
	return string(v)
}
