package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

type Keys struct {
	CLIENT_SECRET        string
	CLIENT_ID            string
	RIOT_API_KEY         string
	TWITCH_CLIENT_ID     string
	TWITCH_CLIENT_SECRET string
	PATCH                string
	SEASON               string
}

var keys Keys

func (a *App) Run() {
	a.InitializeDB()
	a.InitializeRoutes()
	a.InitializeAuth()
	InitializeKeys()

	a.Serve()
}

func InitializeKeys() {
	jsonFile, err := os.Open("keys.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &keys)
}

func (a *App) Serve() {
	http.ListenAndServe(":8000", a.Router)
}
