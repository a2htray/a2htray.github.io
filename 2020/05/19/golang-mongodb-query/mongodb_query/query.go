package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func connect(username string, password string, host string, port int) *mongo.Client {
	credential := options.Credential{
		Username: username,
		Password: password,
	}
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)).SetAuth(credential)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	username := "mongoadmin"
	password := "secret"
	host := "172.28.99.78"
	port := 27017
	database := "mllab"
	client := connect(username, password, host, port)
	collection := client.Database(database).Collection("users")
	ctx := context.TODO()

	var result []bson.D

	cursor, _ := collection.Find(ctx, bson.M{"Username": "booker12"})
	_ = cursor.All(ctx, &result)
	fmt.Println(result)

	cursor, _ = collection.Find(ctx, bson.M{"Identifier": bson.D{{"$gte", 3000}}})
	_ = cursor.All(ctx, &result)
	fmt.Println(len(result))

	cursor, _ = collection.Find(ctx, bson.M{
		"Username": bson.D{bson.E{
			Key: "$in",
			Value: []string{"booker12", "smith79"},
		}},
	})
	_ = cursor.All(ctx, &result)
	fmt.Println(len(result))

	cursor, _ = collection.Find(ctx, bson.M{
		"$or": []interface{}{
			bson.M{"Username": "booker12"},
			bson.M{"Username": "smith79"},
		},
	})
	_ = cursor.All(ctx, &result)
	fmt.Println(len(result))

	cursor, _ = collection.Find(ctx, bson.M{
		"$or": []interface{}{
			bson.M{"Username": "booker12"},
			bson.M{"Username": "smith79"},
		},
	})
	_ = cursor.All(ctx, &result)
	fmt.Println(len(result))
}