package main

import (
	"leafgame/cmd/gateway/center"
	"leafgame/pkg/config"

	"leafgame/pkg/leaf"
)

func main() {
	config.Init()
	leaf.Run(
		center.Module,
	)
}
