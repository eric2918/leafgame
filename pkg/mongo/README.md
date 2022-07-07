## mongodb 计数器
```bash
// 启动数据库初始化
mongo.EnsureCounter(mongo.LOGIN_DB, "counter", "account")
mongo.EnsureCounter(mongo.GAME_DB, "counter", "player")

// 获取返回计数
for i := 0; i < 3; i++ {
    seq, err := mongo.NextSeq(mongo.LOGIN_DB, "counter", "account")
    fmt.Println(i, seq, err)
}
```
