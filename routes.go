package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) InitializeRoutes() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("..")))
	r.HandleFunc("/test", test)

	a.Router = r
}

func test(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "success")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
