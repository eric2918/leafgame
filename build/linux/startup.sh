#!/bin/sh
go mod tidy
start(){
  echo "编译 $1 server"
  go build -o bin/$1 cmd/$1/main.go

  echo "启动 $1 server"
  nohup ./bin/$1 $2 >> /dev/null 2>&1 &
}

start gateway ./bin/conf/gateway.json
start account ./bin/conf/account.json
start game ./bin/conf/game.json
start chat ./bin/conf/chat.json
