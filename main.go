package main

import(
	"log"

	"./database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func errH(message string, err error){ //error handler
	if err != nil { log.Fatalf("[Error] %s: %s",message,err) }
}

func main(){
	//region mongoDB initialization
	db, ctx, err := database.Connect("mongodb://localhost:27017") //TODO: add config file	
	errH("Could not initialize MongoDB client",err)

	defer func() {
		err = db.Disconnect(ctx)
		errH("Could not disconnect from MongoDB client",err)
	}()

	err = db.Ping()
	errH("Could not ping database",err)
	//endregion
	
	result, err := db.CreateTask("Title of the task", "Short description")
	id := result.InsertedID.(primitive.ObjectID).Hex()
	errH("Failed to create a task",err)
	log.Printf("Task no. %s added.",result.InsertedID.(primitive.ObjectID).Hex())
	
	readResult, err := db.ReadTask(id)
	errH("Failed to read a task",err)
	log.Printf("%+v\n",readResult)

	delResult, err := db.DeleteTask(id)
	errH("Failed to remove a task",err)
	log.Printf("%d task removed.",delResult.DeletedCount)

	readResult, err = db.ReadTask(id)
	errH("Failed to read a task",err)
	log.Printf("%+v\n",readResult)
}