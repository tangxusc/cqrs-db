package mongo_impl

import (
	"context"
	"fmt"
	"github.com/siddontang/go/bson"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestNewMongoServer(t *testing.T) {
	//opt := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", `127.0.0.1`, `27018`))
	opt := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", `127.0.0.1`, `3307`))
	opt.SetMaxPoolSize(uint64(1))
	todo := context.TODO()
	client, e := mongo.Connect(todo, opt)
	if e != nil {
		logrus.Errorf("[db]connection mongodb error:%v", e)
		return
	}

	e = Insert(client, todo)
	if e != nil {
		panic(e)
	}

	Find(client, todo, e)

}

func Find(client *mongo.Client, todo context.Context, e error) {
	one := client.Database("aggregate").Collection("a1_aggregate").FindOne(todo,
		bson.M{
			"id": "1",
		})
	fmt.Println("one.Err():", one.Err())
	data := make(map[string]interface{})
	e = one.Decode(&data)
	fmt.Println(data, e)
}

func Insert(client *mongo.Client, todo context.Context) error {
	result, e := client.Database("aggregate").Collection(`a1_event`).InsertOne(todo,
		bson.M{
			`aggId`:     `1`,
			`eventType`: `test1`,
			`data`: bson.M{
				`name`: `test1`,
			},
			`version`:    1,
			`createTime`: `2006-01-02 15:04:05`,
		},
	)
	if e != nil {
		return e
	} else {
		fmt.Println(result.InsertedID)
	}
	return nil
}
