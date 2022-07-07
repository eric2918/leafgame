package client

import (
	"leafgame/msg"
	"leafgame/pb"
	"leafgame/pkg/code"
	"leafgame/pkg/leaf/log"
)

func init() {
	msg.Processor.SetHandler(&pb.AS2C_Login{}, handleLogin)
	msg.Processor.SetHandler(&pb.GS2C_CheckLogin{}, handleCheckLogin)
	msg.Processor.SetHandler(&pb.GS2C_CreatePlayer{}, handleCreatePlayer)

	msg.Processor.SetHandler(&pb.GS2C_GetSkills{}, handleSkills)
	msg.Processor.SetHandler(&pb.GS2C_GetRoles{}, handleRoles)
	msg.Processor.SetHandler(&pb.GS2C_GetUserTeams{}, handleGetUserTeams)
	msg.Processor.SetHandler(&pb.GS2C_GetUserRoles{}, handleGetUserRoles)
	msg.Processor.SetHandler(&pb.GS2C_EditUserTeam{}, handleEditUserTeam)
	msg.Processor.SetHandler(&pb.GS2C_Heartbeat{}, handleHeartbeat)
}

func handleLogin(args []interface{}) {
	rsp := args[0].(*pb.AS2C_Login)
	if rsp.Code != 200 {
		log.Release("login user fail: %s", code.Text(rsp.Code))
		Close()
		return
	}

	playerData.AccountId = rsp.Data.AccountId
	Start(rsp.Data.GameAddr)

	sendMsg := &pb.C2GS_CheckLogin{Token: rsp.Data.Token}
	Client.WriteMsg(sendMsg)
}

func handleCheckLogin(args []interface{}) {
	rsp := args[0].(*pb.GS2C_CheckLogin)

	if rsp.Data != nil && rsp.Data.AccountId != 0 {
		playerData.AccountId = rsp.Data.AccountId
		playerData.NickName = rsp.Data.NickName

		log.Release("%v(%v) login user success", playerData.NickName, playerData.AccountId)
	} else {
		// playerData.UserName = playerData.Username
		sendMsg := &pb.C2GS_CreatePlayer{NickName: playerData.UserName}
		Client.WriteMsg(sendMsg)
	}
}

func handleCreatePlayer(args []interface{}) {
	rsp := args[0].(*pb.GS2C_CreatePlayer)
	playerData.PlayerId = rsp.Data.PlayerId

	log.Release("%v(%v) login and create user success", playerData.NickName, playerData.PlayerId)
}

func handleSkills(args []interface{}) {
	rsp := args[0].(*pb.GS2C_GetSkills)
	log.Release("%v get skill success", rsp)
}

func handleRoles(args []interface{}) {
	rsp := args[0].(*pb.GS2C_GetRoles)
	log.Release("%v get roles success", rsp)
}

func handleGetUserTeams(args []interface{}) {
	rsp := args[0].(*pb.GS2C_GetUserTeams)
	log.Release("%v get team  success", rsp)
}

func handleGetUserRoles(args []interface{}) {
	rsp := args[0].(*pb.GS2C_GetUserRoles)
	log.Release("%v get user role  success", rsp)

}

func handleEditUserTeam(args []interface{}) {
	rsp := args[0].(*pb.GS2C_EditUserTeam)
	log.Release("%v edit team  success", rsp)

}

func handleHeartbeat(args []interface{}) {
	rsp := args[0].(*pb.GS2C_Heartbeat)
	log.Release("%v heart beat success", rsp)
}
