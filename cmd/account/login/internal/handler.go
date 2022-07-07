package internal

import (
	"fmt"
	"leafgame/conf"
	"leafgame/pb"
	"leafgame/pkg/code"
	"leafgame/pkg/jwt"
	"leafgame/pkg/leaf/cluster"
	"leafgame/pkg/leaf/gate"
	"leafgame/pkg/md5"
	"leafgame/pkg/mongo"
	"leafgame/pkg/rand"
	"leafgame/pkg/redis"
	"leafgame/pkg/snowflake"
	"reflect"
	"time"

	jwt2 "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func handle(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handle(&pb.C2AS_Login{}, handleLogin)
	handle(&pb.C2GS_Heartbeat{}, handHeartbeat)
}

func handHeartbeat(args []interface{}) {
	agent := args[1].(gate.Agent)

	sendMsg := &pb.GS2C_Heartbeat{
		Timestamp: time.Now().Unix(),
	}

	agent.WriteMsg(sendMsg)
}

func handleLogin(args []interface{}) {
	req := args[0].(*pb.C2AS_Login)
	agent := args[1].(gate.Agent)

	sendMsg := &pb.AS2C_Login{}
	sendErrFunc := func(code int64) {
		sendMsg.Code = code
		agent.WriteMsg(sendMsg)
	}

	if req.Username == "" {
		sendErrFunc(code.UserNameEmpty)
		return
	}

	var accountData pb.Account
	err := mongo.Collection(mongo.LOGIN_DB, mongo.ACCOUNT_COLLECTION).
		Find(bson.M{"username": req.Username}).
		One(&accountData)
	if err == mgo.ErrNotFound {
		salt := rand.RandStr(6)
		password := md5.EncodeMd5(req.Password + salt)

		accountData = pb.Account{
			AccountId: snowflake.GenID(),
			Username:  req.Username,
			Password:  password,
			Salt:      salt,
			CreateAt:  time.Now().Unix(),
			Status:    0,
			Region:    req.Region,
		}

		err = mongo.Collection(mongo.LOGIN_DB, mongo.ACCOUNT_COLLECTION).Insert(accountData)
	}

	if err != nil {
		sendErrFunc(code.InternalServerError)
		return
	} else if accountData.Password != md5.EncodeMd5(req.Password+accountData.Salt) {
		sendErrFunc(code.IncorrectUsernameOrPassword)
		return
	}

	results, err := cluster.CallN("gateway", "GetBestGameInfo", accountData.AccountId, req.Region)
	if err != nil {
		sendErrFunc(code.NoGameServerAvailable)
		return
	}
	//gameName := results[0].(string)
	gameAddr := results[1].(string)

	// 生成token
	j := jwt.NewJWT()
	claims := jwt.CustomClaims{
		AccountId: accountData.AccountId,
		Username:  accountData.Username,
		StandardClaims: jwt2.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(conf.Server.JwtTimeout) * time.Second).Unix(),
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		sendErrFunc(code.InternalServerError)
		return
	}

	// 保存到redis
	key := fmt.Sprintf("token_%d", claims.AccountId)
	if err = redis.Client.Set(key, token, time.Duration(conf.Server.JwtTimeout)*time.Second).Err(); err != nil {
		sendErrFunc(code.IncorrectUsernameOrPassword)
		return
	}

	sendMsg.Code = 200
	sendMsg.Data = &pb.AS2C_LoginData{
		AccountId: claims.AccountId,
		GameAddr:  gameAddr,
		Token:     token,
	}
	agent.WriteMsg(sendMsg)
}
