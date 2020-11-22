package database

import (
	"context"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"../models"
)

type Database struct{	
	Client *mongo.Client
}

func Connect(uri string) (Database, context.Context, error) { 
	log.Printf("Starting MongoDB client...") //TODO: maybe add log library? //TODO: Add prefixes to logs
	log.Printf("Setting remote URL as %s...",uri)

	options := options.Client().ApplyURI(uri) //TODO: check for options
	client, err := mongo.NewClient(options)
	if err != nil {return Database{},nil,err}
	
	//TODO: Database encryption
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	
	err = client.Connect(ctx)
	if err != nil {return Database{},nil,err}

	log.Printf("MongoDB client initialized.")
	return Database{Client: client}, ctx, err
}


func (db *Database) Disconnect(ctx context.Context) error {
	log.Printf("Disconnecting from database...")
	err := db.Client.Disconnect(ctx)
	if err != nil { return err }
	log.Printf("Disconnected.")
	return nil
}

func (db *Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	log.Printf("Pinging database...")
	err := db.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	log.Printf("Pong!")

	return nil
}

func (db *Database) CreateTask(title string, description string) (*mongo.InsertOneResult, error) {
	task := models.Task{Title: title, Desc: description}
	collection := db.Client.Database("dashboard-tasks").Collection("tasks") //TODO: read from config

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	//TODO: verificate data 

	res, err := collection.InsertOne(ctx,task)
	if err != nil { return nil,err }

	return res, nil
}

//TODO: Maybe CreateManyTasks() (for importing etc)?

func (db *Database) ReadTask(id string) (*models.Task, error)  {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() 

	collection := db.Client.Database("dashboard-tasks").Collection("tasks")

	task := models.Task{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil { return nil, err }

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil { return nil, err }

	return &task, err
}

func (db *Database) DeleteTask(id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Client.Database("dashboard-tasks").Collection("tasks")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil { return nil,err }

	return collection.DeleteOne(ctx, bson.M{"_id": objectID})
}
