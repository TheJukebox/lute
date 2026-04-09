package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Lute LuteConfig

type LuteConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	MinioHost string `yaml:"minioHost"`
	MinioPort string `yaml:"minioPort"`
	PostgresHost string `yaml:"postgresHost"`
	PostgresPort string `yaml:"postgresPort"`
	LogLevel string `yaml:"logLevel"`
}

func NewLuteConfig(
	host string,
	port string,
	minioHost string,
	minioPort string,
	postgresHost string,
	postgresPort string,
	logLevel string,
) LuteConfig {
	return LuteConfig{
		Host: host,
		Port: port,
MinioHost: minioHost,
		MinioPort: minioPort,
		PostgresHost: postgresHost,
		PostgresPort: postgresPort,
		LogLevel: logLevel,
	}
}

func ReadLuteConfig(path string) error {
	configFile, err := os.ReadFile(path);
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configFile, &Lute)
	if err != nil {
		return err
	}
	return nil
}
