## 消息ID规则

#### 消息ID为 5 位数 

| 1 | 01 | 01 |
| :------ | :------ | :------ |
| 服务划分 | 功能模块 | 具体消息 |

- 服务划分：1 位数进行表示，比如 1 配置服务；2 登录服务；3 游戏服务。
- 功能模块：2 位数进行表示，比如 01 登录功能；02 编队功能。
- 具体消息：2 位数进行表示，比如 01 用户登录；02 登录返回。


## protobuf 编译文件中添加自定义标签
1. 安装工具
```bash
go get github.com/favadi/protoc-go-inject-tag
```

2. proto文件中添加标签(按需要自行修改)
```bash
message example  {
  int64 UserId = 1;     // 用户ID  @gotags: bson:"user_id,omitempty"
  string Username = 2;  // 用户名  @gotags: bson:"user_name,omitempty"
}
```

3. 生成pb.go 文件
```bash
protoc --go_out=. *.proto
```

4. 修改生成的文件
```bash
protoc-go-inject-tag -input="*.pb.go"
```
