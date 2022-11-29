package server

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func InitConfig(f string, config *ServerConfig) {
	_, err := toml.DecodeFile(f, config)
	if err != nil {
		fmt.Printf("配置文件错误 %v\r\n", err.Error())
	}
	InitLogger(config)
}

type ServerConfig struct {
	Stdout string
	Stderr string
	Broker map[string]BrokerConfig
}

type BrokerConfig struct {
	Protocol string
	Name     string
	Addr     string
	Port     string
	Enable   bool
	Stdout   string
	Stderr   string
}
