
.PHONY: proto
proto:
	rm -rf pb/*.pb.go
	protoc --proto_path=proto --go_out=plugins=grpc:pb proto/*.proto
	protoc-go-inject-tag -input="pb/*.pb.go"

servers = gateway account game chat
.PHONY: clear
clear:
	for server in $(servers); do \
  		rm -rf bin/$$server; \
    	rm -rf images/$$server/$$server; \
		echo $$server; \
	done
	rm -rf logs

.PHONY: build
build:
	go build -o ./bin/gateway ./cmd/gateway/main.go
	go build -o ./bin/account ./cmd/account/main.go
	go build -o ./bin/game ./cmd/game/main.go
	go build -o ./bin/chat ./cmd/chat/main.go

.PHONY: logs
logs:
	mkdir -p logs/gateway
	mkdir -p logs/account
	mkdir -p logs/game
	mkdir -p logs/chat

.PHONY: startup
startup:logs
	./build/linux/startup.sh

.PHONY: shutdown
shutdown:
	./build/linux/shutdown.sh


version=latest
gate = gateway
account = account
game = game
chat = chat

username = eric@51chiyun.com
registry = registry.cn-hangzhou.aliyuncs.com
namespace = mist-erosion
password = Hx&890108

# docker镜像
.PHONY : images
images : gateway account game chat

gateway:
	@echo "building image for $(gateway) ..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./images/$(gateway)/$(gateway) ./cmd/$(gateway)
	cd ./images/$(gateway) && docker build -t $(gateway):$(version) .

.PHONY : account
account:
	@echo "building image for $(account) ..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./images/$(account)/$(account) ./cmd/$(account)
	cd ./images/$(account) && docker build -t $(account):$(version) .

.PHONY : game
game:
	@echo "building image for $(game) ..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./images/$(game)/$(game) ./cmd/$(game)/main.go
	cd ./images/$(game) && docker build -t $(game):$(version) .

.PHONY : chat
chat:
	@echo "building image for $(chat) ..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./images/$(chat)/$(chat) ./cmd/$(chat)
	cd ./images/$(chat) && docker build -t $(chat):$(version) .

.PHONY : push
push:
	# 镜像标签
	$ docker tag $(gateway):$(version) $(registry)/$(namespace)/$(gateway):$(version)
	$ docker tag $(account):$(version) $(registry)/$(namespace)/$(account):$(version)
	$ docker tag $(game):$(version) $(registry)/$(namespace)/$(game):$(version)
	$ docker tag $(chat):$(version) $(registry)/$(namespace)/$(chat):$(version)

	# 登录阿里云镜像仓库
	$ docker login --username=$(username) --password="$(password)" $(registry)

	# push到阿里云镜像仓库
	$ docker push $(registry)/$(namespace)/$(gateway):$(version)
	$ docker push $(registry)/$(namespace)/$(account):$(version)
	$ docker push $(registry)/$(namespace)/$(game):$(version)
	$ docker push $(registry)/$(namespace)/$(chat):$(version)

	# 删除本地docker images
	docker rmi $(gateway):$(version)
	docker rmi $(account):$(version)
	docker rmi $(game):$(version)
	docker rmi $(chat):$(version)

.PHONY : rmi
rmi:
	docker rmi $(registry)/$(namespace)/$(gateway):$(version)
	docker rmi $(registry)/$(namespace)/$(account):$(version)
	docker rmi $(registry)/$(namespace)/$(game):$(version)
	docker rmi $(registry)/$(namespace)/$(chat):$(version)

conf = /Users/eric/hxgame/src/server/bin/conf
# docker容器
.PHONY : run
run:
	docker run -itd --name $(gateway)_$(version) -p 31001:31001 -v $(conf)/$(gateway).json:/conf/$(gateway).json $(registry)/$(namespace)/$(gateway):$(version)
	docker run -itd --name $(account)_$(version) -p 30001:30001 -p 20001:20001 -p 20002:20002 -v $(conf)/$(account).json:/conf/$(account).json $(registry)/$(namespace)/$(account):$(version)
	docker run -itd --name $(game)_$(version) -p 21001:21001 -p 21002:21002 -v $(conf)/$(game).json:/conf/$(game).json $(registry)/$(namespace)/$(game):$(version)
	docker run -itd --name $(chat)_$(version) -p 32001:32001 -v $(conf)/$(chat).json:/conf/$(chat).json $(registry)/$(namespace)/$(chat):$(version)

.PHONY : stop
stop:
	docker stop $(gateway)_$(version)
	docker stop $(account)_$(version)
	docker stop $(game)_$(version)
	docker stop $(chat)_$(version)

.PHONY : start
start:
	docker start $(gateway)_$(version)
	docker start $(account)_$(version)
	docker start $(game)_$(version)
	docker start $(chat)_$(version)

.PHONY : restart
restart:
	docker restart $(gateway)_$(version)
	docker restart $(account)_$(version)
	docker restart $(game)_$(version)
	docker restart $(chat)_$(version)

.PHONY : rm
rm:
	docker rm -f $(gateway)_$(version)
	docker rm -f $(account)_$(version)
	docker rm -f $(game)_$(version)
	docker rm -f $(chat)_$(version)
