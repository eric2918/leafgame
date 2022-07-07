package internal

import (
	"errors"
	"fmt"

	"leafgame/pkg/jwt"
	"leafgame/pkg/leaf/cluster"
	"leafgame/pkg/redis"
)

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {
	handleRpc("CheckToken", CheckToken)
}

func CheckToken(args []interface{}) (res []interface{}, err error) {
	tokenId := args[0].(string)

	j := jwt.NewJWT()
	claims, err := j.ParserToken(tokenId)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("token_%d", claims.AccountId)
	token, err := redis.Client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	if token != tokenId {
		return nil, errors.New("invalid token")
	}

	return []interface{}{claims.AccountId, claims.Username}, nil
}
