package internal

import (
	"errors"
)

func init() {
	skeleton.RegisterCommand("kickAccount", "Usage: kickAccount|accountId", kickAccount)
}

func kickAccount(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 1 {
		err = errors.New("args len is less than 1")
		return
	}

	accountId := args[0].(string)
	ChanRPC.Go("KickAccount", accountId)
	return
}
