package internal

import (
	"leafgame/cmd/game/center"
	"leafgame/cmd/game/config"
	"leafgame/cmd/game/player"
	"leafgame/pb"
	"leafgame/pkg/leaf/log"
	"leafgame/pkg/mongo"

	"leafgame/pkg/leaf/cluster"
	"leafgame/pkg/leaf/gate"
)

func onAgentInit(agent gate.Agent) {
	var configs pb.Config

	var skills []*pb.Skill
	if err := mongo.Collection(mongo.GAME_DB, mongo.SKILLS_COLLECTION).Find(nil).All(&skills); err != nil {
		log.Error("get skills error: %#v \n", err.Error())
	}
	configs.Skills = skills

	var roles []*pb.Role
	if err := mongo.Collection(mongo.GAME_DB, mongo.ROLES_COLLECTION).Find(nil).All(&roles); err != nil {
		log.Error("get roles error: %#v \n", err.Error())
	}
	configs.Roles = roles

	agent.SetConfigData(config.New(&configs))
}

func onAgentDestroy(agent gate.Agent) {
	var accountId int64
	if val, ok := agent.PlayerData().(int64); ok {
		accountId = val
	} else if player, ok := agent.PlayerData().(*player.Player); ok {
		accountId = player.Player.AccountId

		cluster.Go("gateway", "AccountOffline", player.Player.AccountId)
		center.ChanRPC.Go("UserOffline", player.Player.PlayerId, agent)
	}

	if accountId > 0 {
		center.ChanRPC.Go("AccountOffline", accountId, agent)
	}
}
