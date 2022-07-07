package mongo

import (
	"leafgame/conf"
	"leafgame/pkg/leaf/db/mongodb"
	"log"

	"gopkg.in/mgo.v2"
)

var Context *mongodb.DialContext

func init() {
	var err error
	Context, err = mongodb.Dial(conf.Server.MongodbAddr, conf.Server.MongodbSessionNum)
	if err != nil {
		log.Fatal("mongodb init is error:", err)
	}

	InitIndex()
}

func Index(db, collection string, key []string) error {
	return Context.EnsureIndex(db, collection, key)
}

func UniqueIndex(db, collection string, key []string) error {
	return Context.EnsureUniqueIndex(db, collection, key)
}

func Collection(db, collection string) *mgo.Collection {
	session := Context.Ref()
	defer Context.UnRef(session)

	return session.DB(db).C(collection)
}

func EnsureCounter(db, collection string, id string) error {
	return Context.EnsureCounter(db, collection, id)
}

func NextSeq(db, collection, id string) (seq int, err error) {
	seq, err = Context.NextSeq(db, collection, id)
	return
}
