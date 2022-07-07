package msg

import (
	"leafgame/pb"
	"leafgame/pkg/leaf/network/protobuf"
)

var (
	Processor = protobuf.NewProcessor()
)

func init() {
	Processor.Register(C2AS_LOGIN, &pb.C2AS_Login{})
	Processor.Register(AS2C_LOGIN, &pb.AS2C_Login{})
	Processor.Register(C2GS_CHECK_LOGIN, &pb.C2GS_CheckLogin{})
	Processor.Register(GS2C_CHECK_LOGIN, &pb.GS2C_CheckLogin{})
	Processor.Register(C2GS_CREATE_PLAYER, &pb.C2GS_CreatePlayer{})
	Processor.Register(GS2C_CREATE_PLAYER, &pb.GS2C_CreatePlayer{})

	Processor.Register(C2GS_GET_SKILLS, &pb.C2GS_GetSkills{})
	Processor.Register(GS2C_GET_SKILLS, &pb.GS2C_GetSkills{})

	Processor.Register(C2GS_GET_ROLES, &pb.C2GS_GetRoles{})
	Processor.Register(GS2C_GET_ROLES, &pb.GS2C_GetRoles{})

	Processor.Register(C2GS_GET_USER_TEAMS, &pb.C2GS_GetUserTeams{})
	Processor.Register(GS2C_GET_USER_TEAMS, &pb.GS2C_GetUserTeams{})
	Processor.Register(C2GS_GET_USER_ROLES, &pb.C2GS_GetUserRoles{})
	Processor.Register(GS2C_GET_USER_ROLES, &pb.GS2C_GetUserRoles{})
	Processor.Register(C2GS_EDIT_USER_TEAM, &pb.C2GS_EditUserTeam{})
	Processor.Register(GS2C_EDIT_USER_TEAM, &pb.GS2C_EditUserTeam{})

	Processor.Register(C2GS_HEARTBEAT, &pb.C2GS_Heartbeat{})
	Processor.Register(GS2C_HEARTBEAT, &pb.GS2C_Heartbeat{})

	Processor.Register(10, &pb.C2GS_Hello{})
	Processor.Register(11, &pb.GS2C_Hello{})

}
