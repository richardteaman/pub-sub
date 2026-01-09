package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	logsKey := fmt.Sprintf("%s,*", routing.GameLogSlug)
	logsCh, _, err := pubsub.DeclareAndBind(
		amqpCon,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		logsKey,
		pubsub.QueueDurable,
	)
	if err != nil {
		return
	}
	defer logsCh.Close()

	fmt.Println("AMQP connection was successful.")
	fmt.Println("Starting Peril server...")

	gamelogic.PrintServerHelp()

loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			fmt.Println("Sending pause message")
			err = pubsub.PublishJSON(
				amqpCh,
				string(routing.ExchangePerilDirect),
				string(routing.PauseKey),
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				log.Printf("could not publish: %v", err)
			}
		case "resume":
			fmt.Println("Sending pause message")
			err = pubsub.PublishJSON(
				amqpCh,
				string(routing.ExchangePerilDirect),
				string(routing.PauseKey),
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				log.Printf("could not publish: %v", err)
			}
		case "quit":
			fmt.Println("Exiting loop because of quit command")
			break loop
		default:
			fmt.Println("unknown command")
		}
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Shutting down Peril server...")

}
