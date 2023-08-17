package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error getting all log entries:", err.Error())
		return nil, err
	}
	var logs []*LogEntry
	for cursor.Next(ctx) {
		var entry LogEntry
		err := cursor.Decode(&entry)
		if err != nil {
			log.Println("Error decoding cursor:", err.Error())
			return nil, err
		}
		logs = append(logs, &entry)
	}
	defer cursor.Close(ctx)
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error when generating id:", err.Error())
		return nil, err
	}
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		log.Println("Error getting log entry:", err.Error())
		return nil, err
	}
	return &entry, nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error when generating id:", err.Error())
		return nil, err
	}
	return collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "name", Value: l.Name},
			primitive.E{Key: "data", Value: l.Data},
			primitive.E{Key: "updated_at", Value: time.Now()}},
		}})
	// {"$set", bson.D{
	// 	{"name", l.Name},
	// 	{"data", l.Data},
	// 	{"updated_at", time.Now()}
	// }}
	//})
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	return collection.Drop(ctx)
}

func (l *LogEntry) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error when generating id:", err.Error())
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		log.Println("Error deleting log entry:", err.Error())
		return err
	}
	return nil
}
