package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpUrl := "amqp://guest:guest@localhost:5672/"

	amqpCon, err := amqp.Dial(amqpUrl)
	if err != nil {
		return
	}
	defer amqpCon.Close()

	amqpCh, err := amqpCon.Channel()
	if err != nil {
		return
	}
	defer amqpCh.Close()
	err = pubsub.PublishJSON(
		amqpCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)

	fmt.Println("AMQP connection was successful.")
	fmt.Println("Starting Peril server...")

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Shutting down Peril server...")

}
