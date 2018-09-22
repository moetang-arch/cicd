package config

type Config struct {
	DockerConfig DockerConfig
}

type DockerConfig struct {
	Addr string
}
