package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	mongoURL = "mongodb://localhost:27017"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omit"`
}

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Image_url string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Time      string             `json:"time,omitempty" bson:"time,omitempty"`
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotAcceptable)
	w.Write([]byte("{Message : 'Wrong URL/Method'}"))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorHandler(w, r)
		return
	}
	w.Write([]byte("Welcome"))
}
func createUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.URL.Path != "/users" {
		errorHandler(w, r)
		return
	}
	w.Header().Set("content-type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	collection := client.Database("appointy").Collection("users")
	collection.InsertOne(context.Background(), user)
	json.NewEncoder(w).Encode(user)

}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[7:]
	collection := client.Database("appointy").Collection("users")
	userId, _ := primitive.ObjectIDFromHex(id)
	var user User
	collection.FindOne(context.Background(), bson.D{{"_id", userId}}).Decode(&user)
	json.NewEncoder(w).Encode(user)
}

func createPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.URL.Path != "/posts" {
		errorHandler(w, r)
		return
	}
	w.Header().Set("content-type", "application/json")
	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	collection := client.Database("appointy").Collection("posts")
	collection.InsertOne(context.Background(), post)
	json.NewEncoder(w).Encode(post)
}

func getPostsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[7:]
	collection := client.Database("appointy").Collection("posts")
	postId, _ := primitive.ObjectIDFromHex(id)
	var post Post
	collection.FindOne(context.Background(), bson.D{{"_id", postId}}).Decode(&post)
	json.NewEncoder(w).Encode(post)
}

func getPostsByUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Specific post Info"))
}

func connectMongo() {
	client, _ = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Mongodb Database")
}

func main() {
	connectMongo()
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/users", createUsersHandler)
	mux.HandleFunc("/users/", getUsersHandler)
	mux.HandleFunc("/posts", createPostsHandler)
	mux.HandleFunc("/posts/", getPostsHandler)
	mux.HandleFunc("/posts/users/", getPostsByUserHandler)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
