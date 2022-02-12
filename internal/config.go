package internal

type ServerConfig struct {
}

func LoadServerConfig() *ServerConfig {
	return new(ServerConfig)
}
