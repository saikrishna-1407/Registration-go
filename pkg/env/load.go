package env

import "github.com/kelseyhightower/envconfig"

func Load(v interface{}) error {
	return envconfig.Process(string(Get()), v)
}
