package main

import(
	"log"

	"./database"
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
	
	db.CreateTask("Title of the task", "Short description")
}