package env

type Environment string

const (
	Test Environment = "test"
)

const DefaultEnvironment = Test

var current Environment = DefaultEnvironment

func Get() Environment {
	if current == "" {
		return DefaultEnvironment
	}
	return current
}

func IsProd() bool {
	return Get() == "production"
}
