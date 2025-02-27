package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// create a slice that can hold multiple User values.
var users []User

// handler function -->2 parameters->write the response back to the client that made the request(like status of the code ok,created,Bad request etc) ,
// r*http request contains all the information about the incoming http request like request body,request headers
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	//need to pass the json request body
	err := json.NewDecoder(r.Body).Decode(&newUser)
	//checking the json data whether everything is clear or any data is missing if that is the case below code will execute
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request
		return
	}

	// Checking for required fields
	if newUser.Firstname == "" || newUser.Lastname == "" || newUser.Email == "" || newUser.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Checking if email already exists
	for _, user := range users {
		if user.Email == newUser.Email {
			w.WriteHeader(http.StatusConflict) // 409 Conflict
			json.NewEncoder(w).Encode(map[string]string{"error": "User with this email already exists."})
			return
		}
	}

	users = append(users, newUser)

	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully."})
}

func main() {
	r := mux.NewRouter()                              //Create a "traffic controller" for web requests
	r.HandleFunc("/user", createUser).Methods("POST") //Tell the "controller" what to do for specific requests

	fmt.Println("Server listening on port 8080") //just a message to the programmer that the server is starting
	log.Fatal(http.ListenAndServe(":8080", r))   //Start the server and use the traffic controller r to decide what to do with incoming requests.
}
