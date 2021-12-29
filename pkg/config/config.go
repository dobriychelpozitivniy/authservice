package config

import (
	"github.com/spf13/viper"
	"path"
)

type Config struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Key             string `mapstructure:"key"`
	AccessTokenTTL  int    `mapstructure:"access_token_ttl"`
	RefreshTokenTTL int    `mapstructure:"refresh_token_ttl"`
	MongoConfig     `mapstructure:"mongo_config"`
}

type MongoConfig struct {
	DBHost       string `mapstructure:"db_host"`
	DBUsername   string `mapstructure:"db_username"`
	DBPassword   string `mapstructure:"db_password"`
	DBPort       string `mapstructure:"db_port"`
	DBTimeout    int    `mapstructure:"db_timeout"`
	DBName       string `mapstructure:"db_name"`
	DBCollection string `mapstructure:"db_collection"`
}

func Load(cfgPath string) (*Config, error) {
	var config Config

	viper.AddConfigPath(path.Dir(cfgPath))
	viper.SetConfigName(path.Base(cfgPath))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
