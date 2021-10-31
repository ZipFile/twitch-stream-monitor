package config

type Loader interface {
	Load() (*Config, error)
}
