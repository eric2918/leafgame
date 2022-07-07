package rsa

import (
	"fmt"
	"testing"

	"leafgame/pkg/leaf/log"
)

func TestRsa(t *testing.T) {
	//if err := GenerateKey(2048); err != nil {
	//	log.Error("generate key error:%s", err.Error())
	//	return
	//}

	msg := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiNjI5ODgxMDMzNWEyM2U1NjJhZmQ2ZDUyIiwiZ2FtZV9uYW1lIjoiZ2FtZTEiLCJnYW1lX2FkZHIiOiJsb2NhbGhvc3Q6MjEwMDEiLCJleHAiOjE2NTQ2MDMwOTh9"
	fmt.Println(len([]byte(msg)))
	data, err := Encrypt([]byte(msg), "../../public.crt")
	if err != nil {
		log.Error("rsa encrypt error:%s", err.Error())
		return
	}
	fmt.Println(data)

	res, err := Decrypt(data, "../../private.pem")
	if err != nil {
		log.Error("rsa decrypt error:%s", err.Error())
		return
	}
	fmt.Println(string(res))
}
