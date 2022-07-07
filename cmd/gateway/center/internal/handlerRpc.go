package internal

import (
	"errors"
	"leafgame/pkg/leaf/cluster"
	"leafgame/pkg/leaf/log"
	"math"
	"math/rand"
)

var (
	gameInfoMap    = map[string][]*GameInfo{}
	chatInfoMap    = map[string][]*ChatInfo{}
	roomInfoMap    = map[string]*RoomInfo{}
	accountGameMap = map[int64]*GameInfo{}
)

type GameInfo struct {
	serverName     string
	clientCount    int
	maxClientCount int
	clientAddr     string
}

type ChatInfo struct {
	serverName  string
	clientCount int
	clusterAddr string
}

type RoomInfo struct {
	serverName string
}

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	cluster.AgentChanRPC = ChanRPC

	skeleton.RegisterChanRPC("NewServerAgent", NewServerAgent)
	skeleton.RegisterChanRPC("CloseServerAgent", CloseServerAgent)

	handleRpc("GetBestGameInfo", GetBestGameInfo)
	handleRpc("UpdateGameInfo", UpdateGameInfo)
	handleRpc("UpdateChatInfo", UpdateChatInfo)
	handleRpc("GetRoomInfo", GetRoomInfo)
	handleRpc("DestroyRoom", DestroyRoom)
	handleRpc("AccountOffline", AccountOffline)
}

func NewServerAgent(args []interface{}) {
	serverName := args[0].(string)
	agent := args[1].(*cluster.Agent)
	if serverName[:4] == "game" {
		results, err := agent.CallN("GetGameInfo")
		if err == nil {
			clientCount := results[0].(int)
			maxClientCount := results[1].(int)
			clientAddr := results[2].(string)
			region := results[3].(string)
			gameInfo := &GameInfo{
				serverName:     serverName,
				clientCount:    clientCount,
				maxClientCount: maxClientCount,
				clientAddr:     clientAddr,
			}
			gameInfoMap[region] = append(gameInfoMap[region], gameInfo)

			if len(chatInfoMap[region]) > 0 {
				chatInfo := chatInfoMap[region][rand.Intn(len(chatInfoMap[region]))]

				agent.Go("AddClusterClient", map[string]string{chatInfo.serverName: chatInfo.clusterAddr})
			}
		} else {
			log.Error("GetGameInfo is error: %v", err)
		}
	} else if serverName[:4] == "chat" {
		results, err := agent.CallN("GetChatInfo")
		if err == nil {
			clientCount := results[0].(int)
			clusterAddr := results[1].(string)
			region := results[2].(string)

			chatInfo := &ChatInfo{
				serverName:  serverName,
				clientCount: clientCount,
				clusterAddr: clusterAddr,
			}
			chatInfoMap[region] = append(chatInfoMap[region], chatInfo)

			cluster.Broadcast("game", "AddClusterClient", map[string]string{serverName: clusterAddr})
		} else {
			log.Error("GetChatInfo is error: %v", err)
		}
	}
}

func CloseServerAgent(args []interface{}) {
	serverName := args[0].(string)
	if serverName[:4] == "game" {
		_, ok := gameInfoMap[serverName]
		if ok {
			delete(gameInfoMap, serverName)
		}
	} else if serverName[:4] == "chat" {
		_, ok := chatInfoMap[serverName]
		if ok {
			delete(chatInfoMap, serverName)

			cluster.Broadcast("game", "RemoveClusterClient", serverName)
		}
	}
}

func GetBestGameInfo(args []interface{}) ([]interface{}, error) {
	accountId := args[0].(int64)
	region := args[1].(string)

	var ok bool
	var gameInfo *GameInfo
	if gameInfo, ok = accountGameMap[accountId]; !ok {
		//minClientCount := math.MaxInt32
		//for _, _gameInfo := range gameInfoMap {
		//	if _gameInfo.clientCount < minClientCount && _gameInfo.clientCount < _gameInfo.maxClientCount {
		//		gameInfo = _gameInfo
		//	}
		//}

		// 随机
		if counts := len(gameInfoMap[region]); counts > 0 {
			gameInfo = gameInfoMap[region][rand.Intn(counts)]
		}
	}

	if gameInfo == nil {
		return []interface{}{}, errors.New("no game server to alloc")
	} else {
		accountGameMap[accountId] = gameInfo
		log.Debug("%v account ask game info", accountId)
		return []interface{}{gameInfo.serverName, gameInfo.clientAddr}, nil
	}
}

func UpdateGameInfo(args []interface{}) {
	serverName := args[0].(string)
	clientCount := args[1].(int)
	region := args[2].(string)

	for i, info := range gameInfoMap[region] {
		if info.serverName == serverName {
			gameInfoMap[region][i].clientCount = clientCount
			log.Debug("%v server of client count is %v", serverName, clientCount)
			break
		}
	}
}

func UpdateChatInfo(args []interface{}) {
	serverName := args[0].(string)
	clientCount := args[1].(int)
	region := args[2].(string)

	for i, info := range chatInfoMap[region] {
		if info.serverName == serverName {
			chatInfoMap[region][i].clientCount = clientCount
			log.Debug("%v server of client count is %v", serverName, clientCount)
			break
		}
	}
}

func GetRoomInfo(args []interface{}) (serverName interface{}, err error) {
	roomName := args[0].(string)
	region := args[1].(string)
	roomInfo, ok := roomInfoMap[roomName]
	if ok {
		serverName = roomInfo.serverName
	} else {
		if _, ok = chatInfoMap[region]; !ok {
			err = errors.New("no chat server to alloc")
		}

		var chatInfo *ChatInfo
		minClientCount := math.MaxInt32
		for _, _chatInfo := range chatInfoMap[region] {
			if _chatInfo.clientCount < minClientCount {
				chatInfo = _chatInfo
			}
		}

		if chatInfo == nil {
			err = errors.New("no chat server to alloc")
		} else {
			serverName = chatInfo.serverName
			roomInfoMap[roomName] = &RoomInfo{serverName: chatInfo.serverName}
		}
	}
	return
}

func DestroyRoom(args []interface{}) {
	roomNames := args[0].([]string)
	for _, roomName := range roomNames {
		if _, ok := roomInfoMap[roomName]; ok {
			delete(roomInfoMap, roomName)

		}
	}
	log.Debug("%v rooms is destroy", roomNames)
}

func AccountOffline(args []interface{}) {
	accountId := args[0].(int64)
	if _, ok := accountGameMap[accountId]; ok {
		delete(accountGameMap, accountId)
		log.Debug("%v account is offline", accountId)
	}
}
