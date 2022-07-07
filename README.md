# leaf_game
## 服务架构
1. gateway:网关服，负责服务器的管理、动态加载、负载均衡等中心功能；<br>
2. loginServer:登录服，负责账号登录和验证；<br>
3. gameServer:游戏服，负责玩家数据处理和各种业务功能；<br>
4. chatServer:聊天服，负责聊天、房间相关业务；<br>
		
## 启动要求
1. gateway、loginServer 可以通过nginx实现负载均衡，实现服务器的横向扩展；
2. chatServer和frontServer 通过Gateway实现负载均衡，通过加载不同的配置文件，实现服务器的横向扩展；

## 通信协议
支持websocket/tcp + protobuf/json，本服务使用websocket+protobuf

使用 TCP 协议时（len 默认为两个字节，len 和 id 默认使用网络字节序），在网络中传输的消息格式如下：
```bash
-------------------------------
| len | id | protobuf message |
-------------------------------
```

使用 WebSocket 协议时，发送的消息格式如下：
```bash
-------------------------
| id | protobuf message |
-------------------------
```

## 理消息 ID
```bash
0: msg.C2L_Login
1: msg.L2C_Login
2: msg.C2GS_CheckLogin
3: msg.GS2C_CheckLogin
4: msg.C2GS_CreateUser
5: msg.GS2C_CreateUser
6: msg.ChatMsg
7: msg.C2GS_EnterRoom
8: msg.GS2C_EnterRoom
9: msg.C2GS_LeaveRoom
10: msg.GS2C_LeaveRoom
11: msg.C2GS_SendMsg
12: msg.GS2C_SendMsg
13: msg.GS2C_MsgList
```
