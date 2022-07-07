package client

import (
	"leafgame/conf"
	"sync/atomic"

	"leafgame/pkg/leaf/cluster"
)

var (
	clientCount int32
)

func GetClientCount() int {
	return int(atomic.LoadInt32(&clientCount))
}

func AddClientCount(delta int32) {
	count := atomic.AddInt32(&clientCount, delta)
	cluster.Go("gateway", "UpdateChatInfo", conf.Server.ServerName, int(count), conf.Server.Region)
}
