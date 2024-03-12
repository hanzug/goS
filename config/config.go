package config

var Conf *Config

type Config struct {
	Server *Server
}

type Server struct {
	Port      string `yaml:"port"`
	Version   string `yaml:"port"`
	JwtSecret string `yaml:"jwtSecret"`
}
