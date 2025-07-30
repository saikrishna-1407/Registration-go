package config

type Config struct {
	DatabaseURI string `envconfig:"DATABASE_URL" required:"true"`
}
