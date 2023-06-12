package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Name    string
	Version string
	Host    string
	Port    string
}

type PosgreSQLConfig struct {
	Username           string
	Password           string
	Schema             string
	Host               string
	Port               int
	MaxIdleConnections int
	MaxOpenConnections int
	DebugMode          bool
}

type Database struct {
	PostgreSQL PosgreSQLConfig
}

func NewConfig(path string) *Config {
	fmt.Println("Try new config ....")
	viper.SetConfigFile(path + "../config.json")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	conf := Config{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	configuration, _ := json.Marshal(conf)
	fmt.Println(configuration)
	return &conf
}
