package snowflake

import (
	"leafgame/conf"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	startTime := conf.Server.StartTime
	machineID := conf.Server.MachineID
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		logrus.Fatal("snowflake parse time error", err)
	}

	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		logrus.Fatal("snowflake new node error", err)
	}
}

func GenID() int64 {
	return node.Generate().Int64()
}
