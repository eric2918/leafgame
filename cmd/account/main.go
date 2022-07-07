package main

import (
	"leafgame/cmd/account/gate"
	"leafgame/cmd/account/login"
	"leafgame/pkg/config"
	"leafgame/pkg/leaf"
)

func main() {
	config.Init()
	leaf.Run(
		gate.Module,
		login.Module,
	)

}
