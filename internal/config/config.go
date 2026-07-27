package config

const ThunderPath = `C:\Thunder`

type Config struct{}

func Load() (*Config, error) {
	return &Config{}, nil
}

func (c *Config) RepositoryURL() string {
	return "https://repositorio.brsistemas.com.br/Release/"
}
