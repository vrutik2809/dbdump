package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/vrutik2809/dbdump/cmd/root"
)

func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return client, ctx, cancel, err
}

func ping(client *mongo.Client, ctx context.Context) error {
	err := client.Ping(ctx, nil)

	return err
}

func fetchCollections(db *mongo.Database, ctx context.Context) ([]string, error) {
	collections, err := db.ListCollectionNames(ctx, bson.D{})

	return collections, err
}

func fetchAllDocuments(db *mongo.Database, ctx context.Context, collection string) ([]bson.D, error) {
	cursor, err := db.Collection(collection).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var results []bson.D
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, err
}

func bsonDToJSON(b bson.D) (string, error) {
	// Convert the BSON document to a map
	// var m map[string]interface{}
	// err := bson.Unmarshal(b, &m)
	// if err != nil {
	// 	return "", err
	// }

	// Convert the map to a JSON string
	jsonStr, err := json.Marshal(b)
	if err != nil {
		return "", err
	}

	fmt.Println(string(jsonStr))

	return string(jsonStr), nil
}

func bsonDToJson(b bson.D) (string, error){
	data, err := bson.Marshal(b)
	if err != nil {
		return "",err
	}

	var mp map[string]interface{}
	err = bson.Unmarshal(data,&mp)
	if err != nil {
		return "", err
	}

	res,err := json.Marshal(mp)
	return string(res),err
}

func writeBsonDArrayToFile(bsonDArray []bson.D, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("[")

	first := true

	for _, bsonD := range bsonDArray {
		res, err := bsonDToJson(bsonD)
		if err != nil {
			return err
		}

		if !first {
			file.WriteString(",")
		}
		first = false
		file.WriteString(string(res))
	}
	file.WriteString("]")

	return nil
}

func main() {
	// if err := root.RootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	mongouri := "mongodb://localhost:27017"

	client, ctx, cancel, err := connect(mongouri)

	if err != nil {
		log.Fatal(err)
	}

	defer close(client, ctx, cancel)

	if err := ping(client, ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")

	dbName := "test"

	db := client.Database(dbName)

	collections, err := fetchCollections(db, ctx)
	if err != nil {
		log.Fatal(err)
	}

	// delete dump directory if exists
	os.RemoveAll("dump")
	os.Mkdir("dump", 0777)
	os.Chdir("dump")

	for _, collection := range collections {
		fmt.Println("dumping collection:", collection)

		bsonDArray, err := fetchAllDocuments(db, ctx, collection)
		if err != nil {
			log.Fatal(err)
		}

		if err := writeBsonDArrayToFile(bsonDArray, collection+".json"); err != nil {
			log.Fatal(err)
		}
	}

}
