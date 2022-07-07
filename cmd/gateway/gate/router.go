package gate

import (
	"leafgame/cmd/gateway/center"
	"leafgame/msg"
	"leafgame/pb"
)

func init() {
	msg.Processor.SetRouter(&pb.C2GS_Hello{}, center.ChanRPC)
}
