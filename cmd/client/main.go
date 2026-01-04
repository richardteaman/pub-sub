package main

import (
	"fmt"

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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		return
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)

	ch, _, err := pubsub.DeclareAndBind(
		amqpCon,
		routing.ExchangePerilDirect,
		queueName, routing.PauseKey,
		pubsub.QueueTransient,
	)
	if err != nil {
		return
	}
	defer ch.Close()

	fmt.Println("Starting Peril client...")

	select {}
}
