package main

import (
	"leafgame/cmd/demo/game"
	"leafgame/cmd/demo/gate"
	"leafgame/cmd/demo/login"
	"leafgame/pkg/config"

	"leafgame/pkg/leaf"
)

func main() {
	config.Init()
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
