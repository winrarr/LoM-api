package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) InitializeDB() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://LoM-admin:M4p46iUkuUuMGS@lom.yol35.mongodb.net/LoM?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	a.DB = client.Database("LoM")
}
