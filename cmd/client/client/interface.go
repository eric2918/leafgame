package client

import (
	"errors"
	"leafgame/conf"
	"leafgame/pb"
	"leafgame/pkg/tools"

	"github.com/spf13/cast"
)

var (
	playerData = PlayerData{}
)

type PlayerData struct {
	AccountId int64
	PlayerId  int64
	UserName  string
	NickName  string
}

func init() {
	skeleton.RegisterCommand("login", "login account: input accountName, password", login)
	skeleton.RegisterCommand("getSkills", "get all skill", getAllSkills)
	skeleton.RegisterCommand("getRoles", "get all role", getAllRoles)
	skeleton.RegisterCommand("getUserTeams", "get user team", getUserTeams)
	skeleton.RegisterCommand("editUserTeams", "edit user team", editUserTeam)
	skeleton.RegisterCommand("getUserRoles", "edit user role", getUserRoles)
	skeleton.RegisterCommand("heartbeat", "heart beat", heartbeat)
}

func heartbeat(args []interface{}) (ret interface{}, err error) {
	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	msg := &pb.C2GS_Heartbeat{}
	Client.WriteMsg(msg)
	return
}

func login(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 3 {
		err = errors.New("args len is less than 3")
		return
	}

	username := args[0].(string)
	password := args[1].(string)
	region := args[2].(string)
	playerData.UserName = username

	Start(conf.Server.LoginAddr)

	msg := &pb.C2AS_Login{Username: username, Password: password, Region: region}
	Client.WriteMsg(msg)
	return
}

func getAllSkills(args []interface{}) (ret interface{}, err error) {
	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	msg := &pb.C2GS_GetSkills{}
	Client.WriteMsg(msg)
	return
}

func getAllRoles(args []interface{}) (ret interface{}, err error) {
	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	msg := &pb.C2GS_GetRoles{}
	Client.WriteMsg(msg)
	return
}

func getUserTeams(args []interface{}) (ret interface{}, err error) {
	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	msg := &pb.C2GS_GetUserTeams{}
	Client.WriteMsg(msg)
	return
}

func getUserRoles(args []interface{}) (ret interface{}, err error) {
	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	msg := &pb.C2GS_GetUserRoles{}
	Client.WriteMsg(msg)
	return
}

func editUserTeam(args []interface{}) (ret interface{}, err error) {
	if len(args) < 3 {
		err = errors.New("args len is less than 2")
		return
	}

	teamId := args[0].(string)
	teamName := args[1].(string)
	userRoleIds := args[2].(string)

	if Client == nil {
		err = errors.New("net is disconnect, please input login cmd")
		return
	}

	msg := &pb.C2GS_EditUserTeam{
		TeamName:    teamName,
		UserRoleIds: tools.StringToInt64(userRoleIds),
	}

	if teamId != "nil" {
		msg.TeamId = cast.ToInt64(teamId)
	}

	Client.WriteMsg(msg)
	return
}
