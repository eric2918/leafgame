package player

import "leafgame/pb"

type Player struct {
	Player *pb.Player
}

func New(player *pb.Player) *Player {
	return &Player{Player: player}
}
