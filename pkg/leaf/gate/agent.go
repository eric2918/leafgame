package gate

import (
	"leafgame/pkg/leaf/chanrpc"
	"leafgame/pkg/leaf/module"
	"net"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	Skeleton() *module.Skeleton
	ChanRPC() *chanrpc.Server
	PlayerData() interface{}
	SetPlayerData(data interface{})
	ConfigData() interface{}
	SetConfigData(data interface{})
}
