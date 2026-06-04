package services

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoService(dbName string, defaultCollection string) (*MongoService, error) {

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return nil, errors.New("MONGO_URI is not set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &MongoService{
		client: client,
		db:     db,
	}, nil
}

func (m *MongoService) InsertOne(collectionName string, data interface{}) error {

	collection := m.db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, data)
	return err
}

func (m *MongoService) InsertMany(collectionName string, data []interface{}) error {

	collection := m.db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertMany(ctx, data)
	return err
}

type LogEntry struct {
	Level         string      `bson:"level" json:"level"`
	Message       string      `bson:"message" json:"message"`
	ClientMessage *string     `bson:"client_message,omitempty" json:"client_message,omitempty"`
	Task          interface{} `bson:"task,omitempty" json:"task,omitempty"`
}

func (m *MongoService) InsertLog(
	level string,
	message string,
	clientMessage *string,
	task interface{},
) error {

	logEntry := LogEntry{
		Level:         level,
		Message:       message,
		ClientMessage: clientMessage,
		Task:          task,
	}

	collection := m.db.Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, logEntry)
	return err
}

func (m *MongoService) GetLogs(filter interface{}) ([]map[string]interface{}, error) {

	collection := m.db.Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []map[string]interface{}

	for cur.Next(ctx) {
		var doc map[string]interface{}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		results = append(results, doc)
	}

	return results, nil
}

func (m *MongoService) Close() error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(context.TODO())
}

// // // // // // // // //

var MongoServiceInstance *MongoService

func InitMongo() error {
	svc, err := NewMongoService("kudtax_import_db", "logs")
	if err != nil {
		return err
	}

	MongoServiceInstance = svc
	return nil
}

// services.MongoServiceInstance.InsertLog(
// 	"INFO",
// 	"system started",
// 	nil,
// 	map[string]interface{}{"task": "boot"},
// )
