package database

import (
	"context"
	"time"
	"log"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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