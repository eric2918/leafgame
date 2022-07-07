package gate

import (
	"leafgame/cmd/account/login"
	"leafgame/msg"
	"leafgame/pb"
)

func init() {
	msg.Processor.SetRouter(&pb.C2AS_Login{}, login.ChanRPC)
	msg.Processor.SetRouter(&pb.C2GS_Heartbeat{}, login.ChanRPC)
}
