package main

import (
	"encoding/json"
	"io/ioutil"
	"leafgame/pkg/leaf/log"
	"leafgame/pkg/mongo"
)

func main() {
	DropCollections()
	InitConfig()
	//UpdatePlayer()
}

type Skin struct {
	SkinId int64 `json:"skinId" bson:"skin_id"`
	RoleId int64 `json:"roleId" bson:"role_id"`
}

type Skill struct {
	SkillId   int64  `json:"skillId" bson:"skill_id"`
	SkillName string `json:"skillName" bson:"skill_name"`
}

type Role struct {
	RoleId   int64   `json:"roleId" bson:"role_id"`
	RoleName string  `json:"roleName" bson:"role_name"`
	SkinId   int64   `json:"skinId" bson:"skin_id"`
	SkillIds []int64 `json:"skillIds" bson:"skill_ids"`
	Hp       int64   `json:"hp" bson:"hp"`
	Attack   int64   `json:"attack" bson:"attack"`
	//Element              int     `json:"element" bson:"element"`
	//StoneCtDropRate      int     `json:"stoneCtDropRate" bson:"stone_ct_drop_rate"`
	//StoneSpDropRate      int     `json:"stoneSpDropRate" bson:"stone_sp_drop_rate"`
	//StoneCtCollectWeight int     `json:"stoneCtCollectWeight" bson:"stone_ct_collect_weight"`
	//StoneSpCollectWeight int     `json:"stoneSpCollectWeight" bson:"stone_sp_collect_weight"`
	//IsForbidden          bool    `json:"isForbidden" bson:"is_forbidden"`
}

var path = "./bin/config/"
var configs = []string{"skins", "skills", "roles"}
var roleSkinMap = make(map[int64]int64)

func InitConfig() {
	for _, name := range configs {
		data, err := ioutil.ReadFile(path + name + ".json")
		if err != nil {
			log.Fatal("err: %s", err.Error())
		}

		switch name {
		case "skins":
			var skins []Skin
			err = json.Unmarshal(data, &skins)
			if err != nil {
				log.Fatal("err:%s", err.Error())
			}
			for _, skin := range skins {
				roleSkinMap[skin.RoleId] = skin.SkinId
			}
		case "skills":
			var skills []Skill
			err = json.Unmarshal(data, &skills)
			if err != nil {
				log.Fatal("err:%s", err.Error())
			}

			var docs []interface{}
			for _, r := range skills {
				docs = append(docs, r)
			}

			if err = mongo.Collection(mongo.GAME_DB, name).Insert(docs...); err != nil {
				log.Error("init skill config error: %s", err.Error())
				continue
			}
		case "roles":
			var roles []Role
			err = json.Unmarshal(data, &roles)
			if err != nil {
				log.Fatal("err:", err)
			}

			var docs []interface{}
			for _, r := range roles {
				r.SkinId = roleSkinMap[r.RoleId]
				docs = append(docs, r)
			}

			if err = mongo.Collection(mongo.GAME_DB, name).Insert(docs...); err != nil {
				log.Error("init role config error: %s", err.Error())
				continue
			}
		}
	}
}

func DropCollections() {
	for _, collection := range configs {
		if collection == "skins" {
			continue
		}
		if err := mongo.Collection(mongo.GAME_DB, collection).DropCollection(); err != nil {
			log.Error("drop %s collection error: %s", collection, err.Error())
		}
	}
}

func UpdatePlayer() {

	//info, err := mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).UpdateAll(bson.M{"playerid": 13532180605083655, "userteams.teamid": 13532180605083654}, bson.M{"$set": bson.M{"userteams.$.teamname": "demo"}})
	//if err != nil {
	//	fmt.Println("error:", err.Error())
	//}
	//fmt.Printf("info:%#v \n", info)

	//playerId := 16721836268101640
	//teamId := 16721836268101638
	//team := pb.UserTeam{
	//	TeamId:      111111,
	//	TeamName:    "test",
	//	UserRoleIds: []int64{1, 2, 3, 4},
	//}
	//_ = playerId
	//_ = teamId
	//_ = team

	//	内嵌子集合
	//	更新 $set, 条件中添加自己和查询条件
	//selector := bson.M{"player_id": playerId, "teams.team_id": teamId}
	//update := bson.M{"$set": bson.M{"teams.$": team}}
	//if err := mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).Update(selector, update); err != nil {
	//	fmt.Println("error:", err.Error())
	//}

	//	插入 $push
	//selector := bson.M{"player_id": playerId}
	//update := bson.M{"$push": bson.M{"teams": team}}
	//if err := mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).Update(selector, update); err != nil {
	//	fmt.Println("error:", err.Error())
	//}

	//	删除 $unset
	//selector := bson.M{"player_id": playerId}
	//update := bson.M{"$unset": bson.M{"teams": 1}}
	//if err := mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).Update(selector, update); err != nil {
	//	fmt.Println("error:", err.Error())
	//}

	//selector := bson.M{"player_id": playerId}
	//update := bson.M{"$pop": bson.M{"roles": 1}}
	//if err := mongo.Collection(mongo.GAME_DB, mongo.PLAYER_COLLECTION).Update(selector, update); err != nil {
	//	fmt.Println("error:", err.Error())
	//}

	// 从数组field内删除一个等于value值
	//selector := bson.M{"name": "ls"}
	//update := bson.M{"$pull": bson.M{"tags": "a"}}
	//if err := mongo.Collection(mongo.GAME_DB, "person").Update(selector, update); err != nil {
	//	fmt.Println("error:", err.Error())
	//}

	//	https://www.runoob.com/mongodb/mongodb-atomic-operations.html

}
