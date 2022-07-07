package internal

import (
	"leafgame/pb"
	"reflect"

	"leafgame/pkg/leaf/gate"
)

func handle(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handle(&pb.C2GS_Hello{}, handleHello)
}

func handleHello(args []interface{}) {
	req := args[0].(*pb.C2GS_Hello)
	agent := args[1].(gate.Agent)

	sendMsg := &pb.GS2C_Hello{
		Name: req.Name,
	}
	agent.WriteMsg(sendMsg)
}
