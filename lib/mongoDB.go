package lib

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mapItem struct {
	clinet *mongo.Client
	db     string
}

var mongoDBMap map[string]*mapItem

type MongoDB struct {
	Name     string
	Addr     string
	User     string
	Password string
	DB       string
}

func NewMongoDB(name, addr, user, password string, db string) *MongoDB {
	if mongoDBMap == nil {
		mongoDBMap = make(map[string]*mapItem)
	}
	return &MongoDB{
		Name:     name,
		Addr:     addr,
		User:     user,
		Password: password,
		DB:       db,
	}
}

func (r *MongoDB) InitMongoDB() {
	if mongoDBMap == nil {
		mongoDBMap = make(map[string]*mapItem)
	}
	if r.Name == "" {
		r.Name = "default"
	}
	if _, ok := mongoDBMap[r.Name]; ok {
		return
	}
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s", r.User, r.Password, r.Addr)).SetMaxPoolSize(50) //设置连接池的最大连接数为 50
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("MongoDB connect", err)
		return
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("MongoDB ping", err)
		return
	}
	mongoDBMap[r.Name] = &mapItem{clinet: client, db: r.DB}
}

func GetMongoDB(name ...string) *mongo.Database {
	if mongoDBMap == nil {
		return nil
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := mongoDBMap[n]; ok {
		item := mongoDBMap[n]
		return item.clinet.Database(item.db)
	} else {
		return nil
	}
}

func DestroyMongoDB(name ...string) {
	if mongoDBMap == nil {
		return
	}
	var n string
	if len(name) > 0 {
		n = name[0]
	} else {
		n = "default"
	}
	if _, ok := mongoDBMap[n]; ok {
		_ = mongoDBMap[n].clinet.Disconnect(context.TODO())
		delete(mongoDBMap, n)
	}
}

func DestroyMongoDBAll() {
	if mongoDBMap == nil {
		return
	}
	for k, r := range mongoDBMap {
		_ = r.clinet.Disconnect(context.TODO())
		delete(mongoDBMap, k)
	}
}
