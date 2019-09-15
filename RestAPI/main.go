package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Inbox struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Email    string    `json:"email,omitempty"`
	Messages *Messages `json:"messages,omitempty"`
}
type Messages struct {
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
}

var inbox []Inbox

func getInboxEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(inbox)
}

func getMessageEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range inbox {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Inbox{})
}
func createMessageEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var message Inbox
	_ = json.NewDecoder(req.Body).Decode(&message)
	message.ID = params["id"]
	inbox = append(inbox, message)
	json.NewEncoder(w).Encode(inbox)
}

func deleteMessageEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range inbox {
		if item.ID == params["id"] {
			inbox = append(inbox[:index], inbox[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(inbox)
}

func main() {
	router := mux.NewRouter()
	inbox = append(inbox, Inbox{ID: "1", Name: "Hartanto Setiawan", Email: "hartantosetiawaan@gmail.com", Messages: &Messages{Subject: "Test", Message: "Hello world"}})
	inbox = append(inbox, Inbox{ID: "2", Name: "Alex", Email: "Gart@gmail.com"})
	router.HandleFunc("/message", getInboxEndpoint).Methods("GET")
	router.HandleFunc("/message/{id}", getMessageEndpoint).Methods("GET")
	router.HandleFunc("/message/{id}", createMessageEndpoint).Methods("POST")
	router.HandleFunc("/message/{id}", deleteMessageEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4321", router))
}
