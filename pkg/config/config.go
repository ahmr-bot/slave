package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	Debug struct {
		Debug bool
	}
	Sync struct {
		Endpoint      string
		MinioEndpoint string
		AccessKey     string
		SecretKey     string
		UseSSL        bool
		BucketName    string
	}
}

var conf Config

func LoadConfig(file string) error {
	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return err
	}
	return err
}

func GetConfig() *Config {
	return &conf
}
