package internal

import (
	"fmt"
	"leafgame/cmd/game/center"
	"leafgame/cmd/game/config"
	"leafgame/cmd/game/player"
	"leafgame/conf"
	"leafgame/msg"
	"leafgame/pb"
	"leafgame/pkg/code"
	"leafgame/pkg/leaf/log"
	"leafgame/pkg/mongo"
	"leafgame/pkg/snowflake"
	"time"

	"leafgame/pkg/leaf/cluster"
	"leafgame/pkg/leaf/gate"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	msg.Processor.SetHandler(&pb.C2GS_CheckLogin{}, handleCheckLogin)
	msg.Processor.SetHandler(&pb.C2GS_CreatePlayer{}, handleCreatePlayer)
	msg.Processor.SetHandler(&pb.C2GS_GetSkills{}, handleSkills)
	msg.Processor.SetHandler(&pb.C2GS_GetRoles{}, handleRoles)
	msg.Processor.SetHandler(&pb.C2GS_GetUserTeams{}, handleGetUserTeams)
	msg.Processor.SetHandler(&pb.C2GS_GetUserRoles{}, handleGetUserRoles)
	msg.Processor.SetHandler(&pb.C2GS_EditUserTeam{}, handleEditUserTeam)
	msg.Processor.SetHandler(&pb.C2GS_Heartbeat{}, handHeartbeat)
}

func handHeartbeat(args []interface{}) {
	agent := args[1].(gate.Agent)

	timestamp := time.Now().Unix()

	sendMsg := &pb.GS2C_Heartbeat{
		Timestamp: timestamp,
	}

	// 记录最后心跳时间，定时或推出时更新数据库
	playerData := agent.PlayerData().(*player.Player)
	playerData.Player.LastHeartbeatTime = timestamp
	fmt.Println(playerData.Player)

	agent.WriteMsg(sendMsg)
}

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

func handleCheckLogin(args []interface{}) {
	req := args[0].(*pb.C2GS_CheckLogin)
	agent := args[1].(gate.Agent)

	sendMsg := &pb.GS2C_CheckLogin{
		Code: 200,
		Data: nil,
	}
	res, err := cluster.CallN("account", "CheckToken", req.Token)
	account := map[string]interface{}{
		"AccountId": res[0],
		"Username":  res[1],
	}
	accountId := account["AccountId"].(int64)

	if err != nil {
		sendMsg.Code = code.InternalServerError
		agent.WriteMsg(sendMsg)
		return
	}

	for {
		ok, err := center.ChanRPC.Call1("AccountOnline", accountId, agent)
		if err != nil {
			sendMsg.Code = code.InternalServerError
			agent.WriteMsg(sendMsg)
			return
		}

		if ok.(bool) {
			break
		} else {
			time.Sleep(time.Second)
		}
	}

	var playerInfo pb.Player
	err = mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).Find(bson.M{"account_id": accountId}).One(&playerInfo)
	if err == nil {
		sendMsg.Code = 200
		sendMsg.Data = &pb.GS2C_CheckLoginData{
			AccountId: playerInfo.AccountId,
			NickName:  playerInfo.NickName,
		}

		agent.SetPlayerData(player.New(&playerInfo))
		center.ChanRPC.Go("UserOnline", playerInfo.PlayerId, agent)
	} else if err == mgo.ErrNotFound {
		agent.SetPlayerData(account)
	} else {
		sendMsg.Code = code.InternalServerError
	}
	agent.WriteMsg(sendMsg)
}

func handleCreatePlayer(args []interface{}) {
	req := args[0].(*pb.C2GS_CreatePlayer)
	agent := args[1].(gate.Agent)

	sendMsg := &pb.GS2C_CreatePlayer{}
	account, ok := agent.PlayerData().(map[string]interface{})
	accountId := account["AccountId"].(int64)
	if !ok {
		sendMsg.Code = code.InternalServerError
		agent.WriteMsg(sendMsg)
		return
	}

	if req.NickName == "" {
		req.NickName = account["Username"].(string)
	}

	//var playerInfo pb.Player
	playerInfo, err := InitPlayer(accountId, req.NickName)
	if err != nil {
		sendMsg.Code = code.IncorrectUsernameOrPassword
		agent.WriteMsg(sendMsg)
		return
	}

	if err = mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).Insert(playerInfo); err != nil {
		sendMsg.Code = code.InternalServerError
		agent.WriteMsg(sendMsg)
		return
	} else {
		sendMsg.Code = 200
		sendMsg.Data = &pb.GS2C_CreatePlayerData{
			PlayerId: playerInfo.PlayerId,
		}

		agent.SetPlayerData(player.New(&playerInfo))
		center.ChanRPC.Go("UserOnline", playerInfo.PlayerId, agent)
	}

	agent.WriteMsg(sendMsg)
}

func InitPlayer(accountId int64, nickname string) (info pb.Player, err error) {
	// 初始化用户角色
	var allRoles []pb.Role
	if err = mongo.Collection(mongo.GAME_DB, mongo.ROLES_COLLECTION).Find(nil).All(&allRoles); err != nil {
		return
	}

	var userRoles []*pb.UserRole
	var userRoleIds []int64
	for _, role := range allRoles {
		userRole := &pb.UserRole{
			UserRoleId: snowflake.GenID(),
			RoleId:     role.RoleId,
			Level:      1,
		}
		userRoles = append(userRoles, userRole)
		userRoleIds = append(userRoleIds, userRole.UserRoleId)
	}

	var userTeams []*pb.UserTeam
	for i := 1; i <= conf.Server.MaxTeamsCounts; i++ {
		userTeams = append(userTeams, &pb.UserTeam{
			TeamId: snowflake.GenID(),
		})
	}

	// 初始化用户编队
	info = pb.Player{
		PlayerId:  snowflake.GenID(),
		AccountId: accountId,
		NickName:  nickname,
		Avatar:    "",
		Roles:     userRoles,
		Teams:     userTeams,
	}
	return
}

func handleSkills(args []interface{}) {
	agent := args[1].(gate.Agent)

	sendMsg := &pb.GS2C_GetSkills{
		Code: 200,
		Data: nil,
	}

	configData := agent.ConfigData().(*config.Config)
	sendMsg.Data = configData.Skills

	agent.WriteMsg(sendMsg)
}

func handleRoles(args []interface{}) {
	agent := args[1].(gate.Agent)

	sendMsg := &pb.GS2C_GetRoles{
		Code: 200,
		Data: nil,
	}

	configData := agent.ConfigData().(*config.Config)
	sendMsg.Data = configData.Roles

	agent.WriteMsg(sendMsg)
}

func handleGetUserTeams(args []interface{}) {
	agent := args[1].(gate.Agent)

	playerData := agent.PlayerData().(*player.Player)

	sendMsg := &pb.GS2C_GetUserTeams{
		Code: 200,
		Data: nil,
	}

	sendMsg.Data = playerData.Player.Teams

	agent.WriteMsg(sendMsg)
}

func handleGetUserRoles(args []interface{}) {
	agent := args[1].(gate.Agent)

	playerData := agent.PlayerData().(*player.Player)

	sendMsg := &pb.GS2C_GetUserRoles{
		Code: 200,
		Data: nil,
	}

	sendMsg.Data = playerData.Player.Roles

	agent.WriteMsg(sendMsg)
}

func handleEditUserTeam(args []interface{}) {
	req := args[0].(*pb.C2GS_EditUserTeam)

	agent := args[1].(gate.Agent)
	playerData := agent.PlayerData().(*player.Player)

	sendMsg := &pb.GS2C_EditUserTeam{
		Code: 200,
		Data: nil,
	}

	teamId := req.TeamId
	add := false
	if teamId == 0 {
		teamId = snowflake.GenID()
		add = true
	}

	team := pb.UserTeam{
		TeamId:      teamId,
		TeamName:    req.TeamName,
		UserRoleIds: req.UserRoleIds,
	}
	if add {
		playerData.Player.Teams = append(playerData.Player.Teams, &team)
	} else {
		for i, userTeam := range playerData.Player.Teams {
			if userTeam.TeamId == teamId {
				playerData.Player.Teams[i] = &team
			}
		}
	}

	update := bson.M{"$set": bson.M{"teams": playerData.Player.Teams}}
	selector := bson.M{"player_id": playerData.Player.PlayerId}
	err := mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).Update(selector, update)
	if err != nil {
		log.Error("edit user group error : %v \n", err)
		sendMsg.Code = code.InternalServerError
	}

	// 更新缓存
	agent.SetPlayerData(playerData)

	sendMsg.Data = playerData.Player.Teams

	agent.WriteMsg(sendMsg)
}
