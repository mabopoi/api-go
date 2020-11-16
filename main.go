package main

import (
	"net/http"
	"context"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var client *mongo.Client

func main(){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
	router := mux.NewRouter()
	router.HandleFunc("/tasks", PostEndpoint).Methods("POST")
	router.HandleFunc("/tasks", GetEndpoint).Methods("GET")
	router.HandleFunc("/tasks/{id}", PatchEndpoint).Methods("PATCH")
	router.HandleFunc("/tasks/{id}", DeleteEndpoint).Methods("DELETE")
	http.ListenAndServe(":3001", router)
}