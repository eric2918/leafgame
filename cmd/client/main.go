package main

import (
	client2 "leafgame/cmd/client/client"
	"leafgame/msg"
	"leafgame/pkg/leaf"
)

func main() {
	client2.Init(msg.Processor)

	leaf.Run(
		client2.Module,
	)
}
