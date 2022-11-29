package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/logrusorgru/aurora"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
)

func main() {

	config := mqtt.ServerConfig{}
	configFile := flag.String("f", "./mqtt.toml", "mqtt 配置文件路径 默认mqtt.toml")
	flag.Parse()
	if *configFile != "" {
		mqtt.InitConfig(*configFile, &config)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	fmt.Println(aurora.Magenta("启动 MQTT Broker"))
	server := mqtt.NewServer(nil)
	for _, broker := range config.Broker {
		if broker.Protocol == "tcp" {
			tcp := listeners.NewTCP(broker.Name, broker.Addr)
			err := server.AddListener(tcp, nil)
			if err != nil {
				mqtt.Logger.Sugar().Errorf("%v", err)
			} else {
				mqtt.Logger.Sugar().Infof("监听 tcp %v", broker.Addr)
			}
		}
		if broker.Protocol == "websocket" {
			ws := listeners.NewWebsocket(broker.Name, broker.Addr)
			err := server.AddListener(ws, nil)
			if err != nil {
				mqtt.Logger.Sugar().Errorf("%v", err)
			} else {
				mqtt.Logger.Sugar().Infof("监听 websocket %v", broker.Addr)
			}
		}
	}

	go server.Serve()
	fmt.Println(aurora.BgMagenta("  启动完成  "))
	<-done
	fmt.Println(aurora.BgRed("  捕捉到信号  "))
	server.Close()
	fmt.Println(aurora.BgGreen("  完成  "))
}
