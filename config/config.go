package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Database struct{
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		Name string `yaml:"name"`
	}`yaml:"database"`
}

func NewConfig() *Config {
	cfg := &Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("issue in unmarshal config: %w", err)
	}
	return cfg
}

func (cfg *Config) GetDatabaseName() string {
	return cfg.Database.Name
}