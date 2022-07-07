package main

import (
	"leafgame/cmd/gateway/center"
	"leafgame/cmd/gateway/gate"
	"leafgame/pkg/config"

	"leafgame/pkg/leaf"
)

func main() {
	config.Init()
	leaf.Run(
		gate.Module,
		center.Module,
	)
}
