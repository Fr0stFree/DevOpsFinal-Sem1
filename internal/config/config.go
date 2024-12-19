package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

const DateFormat = "2006-01-02"

type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
}

type DBConfig struct {
	Host     string `yaml:"port"`
	Port     int    `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

func Load(path, fmt string) (Config, error) {
	log.Printf("loading configuration from '%s'\n", path)
	viper.SetConfigFile(path)
	viper.SetConfigType(fmt)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	log.Println("configuration loaded successfully")
	return config, nil
}
