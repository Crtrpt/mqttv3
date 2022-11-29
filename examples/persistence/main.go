package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/logrusorgru/aurora"
	"go.etcd.io/bbolt"

	mqtt "github.com/crtrpt/mqtt/broker"
	"github.com/crtrpt/mqtt/broker/listeners"
	"github.com/crtrpt/mqtt/broker/listeners/auth"
	"github.com/crtrpt/mqtt/broker/persistence/bolt"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	fmt.Println(aurora.Magenta("Mochi MQTT Server initializing..."), aurora.Cyan("Persistence"))

	server := mqtt.NewServer(nil)
	tcp := listeners.NewTCP("t1", ":1883")
	err := server.AddListener(tcp, &listeners.Config{
		Auth: new(auth.Allow),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = server.AddStore(bolt.New("mochi-test.db", &bbolt.Options{
		Timeout: 500 * time.Millisecond,
	}))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println(aurora.BgMagenta("  Started!  "))

	<-done
	fmt.Println(aurora.BgRed("  Caught Signal  "))

	server.Close()
	fmt.Println(aurora.BgGreen("  Finished  "))
}
