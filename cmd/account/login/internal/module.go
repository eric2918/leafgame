package internal

import (
	"leafgame/base"

	"leafgame/pkg/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {
}
