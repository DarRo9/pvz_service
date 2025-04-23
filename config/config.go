package config

import "github.com/spf13/viper"

type Config struct {
	Cities       []string `mapstructure:"cities"`
	ProductTypes []string `mapstructure:"product_types"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
