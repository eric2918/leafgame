package tools

import "gopkg.in/mgo.v2/bson"

func SliceObjectIdToString(slice []bson.ObjectId) (res []string) {
	for _, objectId := range slice {
		res = append(res, objectId.Hex())
	}
	return
}

func SliceStringToObjectId(slice []string) (res []bson.ObjectId) {
	for _, str := range slice {
		res = append(res, bson.ObjectIdHex(str))
	}
	return
}
