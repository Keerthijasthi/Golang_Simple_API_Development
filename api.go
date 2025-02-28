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

func getUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.Email == email && user.Password == password {
			response := map[string]string{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"email":     user.Email,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
	json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password."})
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	for i, user := range users {
		if user.Email == email {
			users = append(users[:i], users[i+1:]...) // Remove the user from the slice
			json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully."})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound) // 404 Not Found
	json.NewEncoder(w).Encode(map[string]string{"error": "User not found."})
}

func main() {
	r := mux.NewRouter()                              //Create a "traffic controller" for web requests
	r.HandleFunc("/user", createUser).Methods("POST") //Tell the "controller" what to do for specific requests
	r.HandleFunc("/user", getUser).Methods("GET")
	r.HandleFunc("/user", deleteUser).Methods("DELETE") //Add delete handler

	fmt.Println("Server listening on port 8080") //just a message to the programmer that the server is starting
	log.Fatal(http.ListenAndServe(":8080", r))   //Start the server and use the traffic controller r to decide what to do with incoming requests.
}
