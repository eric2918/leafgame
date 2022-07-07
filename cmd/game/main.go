package main

import (
	"leafgame/cmd/game/center"
	"leafgame/cmd/game/gate"
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
