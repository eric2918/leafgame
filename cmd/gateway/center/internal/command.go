package internal

import (
	"fmt"

	"leafgame/pkg/leaf/cluster"
)

func init() {
	skeleton.RegisterCommand("echo", "echo account inputs", echo)
	skeleton.RegisterCommand("getGameInfo", "return all game info", getGameInfo)
	skeleton.RegisterCommand("getChatInfo", "return all chat info", getChatInfo)
	skeleton.RegisterCommand("updateConfig", "update game config", updateConfig)
}

func echo(args []interface{}) (ret interface{}, err error) {
	return fmt.Sprintf("%v", args), nil
}

func getGameInfo(args []interface{}) (ret interface{}, err error) {
	ret = fmt.Sprintf("%s", gameInfoMap)
	return
}

func getChatInfo(args []interface{}) (ret interface{}, err error) {
	ret = fmt.Sprintf("%s", chatInfoMap)
	return
}

func updateConfig(args []interface{}) (ret interface{}, err error) {
	// 获取当前在线的游戏服
	// 远程调用游戏服更新配置

	ret1 := " \nupdate config \n"
	for _, gameInfo := range gameInfoMap {
		for _, info := range gameInfo {
			cluster.Go(info.serverName, "UpdateConfig")
			ret1 += fmt.Sprintf("%s %s update success \n", info.serverName, info.clientAddr)
		}
	}
	ret = ret1
	return
}
