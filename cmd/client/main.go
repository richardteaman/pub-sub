package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

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

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		amqpCon,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.QueueTransient,
		handlerPause(gameState),
	)

	fmt.Println("Starting Peril client...")

loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "spawn":
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			_, err = gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Move successful")
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break loop
		default:
			fmt.Println("Unknown command")
		}

	}

}
