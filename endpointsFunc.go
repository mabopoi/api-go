package main

import (
	"net/http"
	"encoding/json"
	"context"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

//Task model 
type Task struct{
	ID 		primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title 	string `json:"title,omitempty" bson:"title,omitempty"` 
	Content string `json:"content,omitempty" bson:"content,omitempty"`
}


//PostEndpoint allows us to create a task 
func PostEndpoint(res http.ResponseWriter, req *http.Request){
	res.Header().Add("content-type", "application/json")
	var task Task 
	json.NewDecoder(req.Body).Decode(&task)
	collection := client.Database("goTest").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, task)
	json.NewEncoder(res).Encode(result)
}

//GetEndpoint allows us to get all tasks 
func GetEndpoint(res http.ResponseWriter, req *http.Request){
	res.Header().Add("content-type", "application/json")
	var tasks []Task
	collection := client.Database("goTest").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var task Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}
	if err := cursor.Err(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(res).Encode(tasks)
}

//PatchEndpoint allows us to change a task
func PatchEndpoint(res http.ResponseWriter, req *http.Request){
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var task Task
	json.NewDecoder(req.Body).Decode(&task)
	var newTask Task
	collection := client.Database("goTest").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOneAndUpdate(ctx,Task{ID: id }, task).Decode(&newTask)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(res).Encode(newTask)
}

//DeleteEndpoint allows us to delete a task 
func DeleteEndpoint(res http.ResponseWriter, req *http.Request){
	res.Header().Add("content-type", "application/json")
	params := mux.Vars(req)
	var task Task
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := client.Database("goTest").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOneAndDelete(ctx, Task{ID: id }).Decode(&task)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(res).Encode(task)
}