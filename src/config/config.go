package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Charset  string `yaml:"charset"`
	} `yaml:"database"`
	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

var AppConfig *Config

func init() {
	config := &Config{}
	file, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	if config.JWT.Secret == "" {
		config.JWT.Secret = "a_very_secret_key_that_should_be_in_config"
		log.Println("JWT secret not found in config, using default. Please set a secret in your config.yaml")
	}

	AppConfig = config
	log.Println("Configuration loaded successfully.")
}
