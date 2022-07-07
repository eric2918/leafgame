package main

import (
	"leafgame/cmd/chat/center"
	"leafgame/cmd/chat/room"
	"leafgame/pkg/config"

	"leafgame/pkg/leaf"
	"leafgame/pkg/leaf/module"
)

func main() {
	config.Init()
	modules := []module.Module{center.Module}
	modules = append(modules, room.CreateModules()...)
	leaf.Run(modules...)
}
