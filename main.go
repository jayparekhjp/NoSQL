package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)

type User struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `json:"name"`
	Password  string        `json:"password"`
	TimeStamp time.Time     `json:"timestamp"`
}

type UserResource struct {
	UsersResource User `json:"user"`
}

type UsersResource struct {
	Users []User `json:"users"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API is ALIVE!")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {

	var userResource UserResource

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	user := userResource.User
	// Generate a new ID
	obj_id := bson.NewObjectId()
	// Transfer this ID to user
	user.ID = obj_id
	// Give TimeStamp to user
	user.TimeStamp = time.Now()
	// Insert user into collection
	err = collection.Insert(&user)
	if err != nil {
		panic(err)
	} else {
		log.Printf("Successfully created User: %s", user.Name)
	}
	j, err := json.Marshal(UserResource{User: user})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {

	var users UsersResource

	// Iteration over collection to get users
	iter := collection.Find(nil).Iter()
	result := User{}
	for iter.Next(&result) {
		users = append(users, result)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(UsersResource{Users: users})
	if err != nil {
		panic(err)
	}
	w.Write(j)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	// Getting UserID from the incoming URL
	vars := mux.Vars(r)
	id := bson.ObjectIdHex(vars["id"])

	// Decode incoming json data
	var userResource UserResource
	err = json.NewDecoder(r.Body).Decode(&userResource)
	// if err != {
	// 	panic(err)
	// }

	// Partial update the existing user
	err = collection.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"name": userResource.User.Name, "password": userResource.User.Password}})
	if err == nil {
		log.Printf("Updated User: %s", id, userResource.User.Name)
	} else {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id := vars["id"]

	// Remove user from database
	err = collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Printf("Could not find the User")
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	// Configuring API endpoints
	r := mux.NewRouter()
	r.HandleFunc("api/ping", pingHandler).Methods("GET")
	r.HandleFunc("api/users", getUsersHandler).Methods("GET")
	r.HandleFunc("api/users/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("api/users", createUserHandler).Methods("POST")
	r.HandleFunc("api/users/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("api/users/{id}", deleteUserHandler).Methods("DELETE")
	http.Handle("/api/", r)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Starting MongoDB Session...")
	var err error
	session, err = mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer sessio.Close()
	sessio.SetMode(mgo.Monotonic, true)
	collection = session.DB("counterBurger").C("users")

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
