package internal

import (
	"errors"
	"leafgame/cmd/chat/client"
	"leafgame/cmd/chat/room"
	"leafgame/conf"
	"leafgame/pkg/leaf/chanrpc"
	"leafgame/pkg/leaf/cluster"
)

var (
	roomModuleMap = map[string]room.Module{}
	routeMoreMap  = map[interface{}]interface{}{}
)

func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

func init() {
	handleRpc("GetChatInfo", GetChatInfo, chanrpc.FuncCommon)
	handleRpc("EnterRoom", EnterRoom, chanrpc.FuncExtRet)
	handleRpc("LeaveRoom", RouteSingle, chanrpc.FuncRoute)
	handleRpc("SendMsg", RouteSingle, chanrpc.FuncRoute)

	routeMoreMap["LeaveRooms"] = "LeaveRoom"
	handleRpc("LeaveRooms", RouteMore, chanrpc.FuncRoute)
}

func GetChatInfo(args []interface{}) ([]interface{}, error) {
	return []interface{}{client.GetClientCount(), conf.Server.ListenAddr, conf.Server.Region}, nil
}

func EnterRoom(args []interface{}) {
	roomName := args[0].(string)

	module := roomModuleMap[roomName]
	if module == nil {
		module = room.GetBestModule()
		if module == nil {
			retFunc := args[len(args)-1].(chanrpc.ExtRetFunc)
			retFunc(nil, errors.New("get best room module rpc is fail"))
			return
		}

		roomModuleMap[roomName] = module
	}

	newArgs := []interface{}{module}
	newArgs = append(newArgs, args...)
	skeleton.AsynCall(module.GetChanRPC(), "EnterRoom", newArgs...)
}

func RouteSingle(args []interface{}) {
	id := args[0]
	roomName := args[1].(string)

	module := roomModuleMap[roomName]
	if module == nil {
		retFunc := args[len(args)-1].(chanrpc.ExtRetFunc)
		retFunc(nil, errors.New("this room is not exist"))
		return
	}

	args = append([]interface{}{module}, args[1:]...)
	skeleton.AsynCall(module.GetChanRPC(), id, args...)
}

func RouteMore(args []interface{}) {
	id := args[0]
	roomNames := args[1].([]string)
	retFunc := args[len(args)-1].(chanrpc.ExtRetFunc)

	id = routeMoreMap[id]
	for _, roomName := range roomNames {
		module := roomModuleMap[roomName]
		if module != nil {
			args = append([]interface{}{module, roomName}, args[2:]...)
			module.GetChanRPC().Go(id, args...)
		}
	}

	retFunc(nil, nil)
}
