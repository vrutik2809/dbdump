package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	username string
	password string
	host string
	port uint
	dbName string
	isSRV bool
	client *mongo.Client
	ctx context.Context
	cancel context.CancelFunc
}

func NewMongoDB(username string, password string, host string, port uint, dbName string,isSRV bool) *mongoDB {
	return &mongoDB{
		username: username,
		password: password,
		host: host,
		port: port,
		dbName: dbName,
		isSRV: isSRV,
		client: nil,
		ctx: nil,
		cancel: nil,
	}
}

func (m *mongoDB) GetURI() string {
	var port string
	if m.port == 0 {
		port = ""
	} else {
		port = fmt.Sprintf(":%d", m.port)
	}

	var prefix string
	if m.isSRV {
		prefix = "mongodb+srv://"
	} else {
		prefix = "mongodb://"
	}

	if m.username == "" || m.password == "" {
		return prefix + m.host + port
	}

	return prefix + m.username + ":" + m.password + "@" + m.host + port
}

func (m *mongoDB) Connect() error {
	uri := m.GetURI()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		defer cancel()
		return err
	}
	m.client = client
	m.ctx = ctx
	m.cancel = cancel
	return nil
}

func (m *mongoDB) Close() {
	defer m.cancel()

	defer func() {
		if err := m.client.Disconnect(m.ctx); err != nil {
			log.Fatal(err)
		}
	}()
}

func (m *mongoDB) Ping() error {
	err := m.client.Ping(m.ctx, nil)

	return err
}

func (m *mongoDB) FetchCollections() ([]string, error) {
	collections, err := m.client.Database(m.dbName).ListCollectionNames(m.ctx, bson.D{})

	return collections, err
}

func (m *mongoDB) FetchAllDocuments(collection string) ([]bson.D, error) {
	cursor, err := m.client.Database(m.dbName).Collection(collection).Find(m.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var results []bson.D
	if err := cursor.All(m.ctx, &results); err != nil {
		return nil, err
	}

	return results, err
}