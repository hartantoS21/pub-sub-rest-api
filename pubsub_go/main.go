package main

import (
	"./pubsub"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoId()  (string){
	return uuid.Must(uuid.NewV4()).String()
}

var ps = & pubsub.PubSub{}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("New Client si connected")
	client :=pubsub.Client{
		Id:autoId(),
		Connection:conn,
	}
	ps.AddClient(client)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("something went wrong ",err)
			ps.RemoveClient(client)
			log.Println("total client and subscription", len(ps.Clients), len(ps.Subcription))
			return
		}
		ps.HandleReceiveMessage(client,messageType,p)
	}
}

func main() {
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		//payload:=map[string]interface{}{
		//	"message":"hello Go",
		//}
		//w.Header().Set("Content-Type","application/json")
		//json.NewEncoder(w).Encode(payload)
		http.ServeFile(w,r, "static")
	})
	http.HandleFunc("/ws",websocketHandler)
	http.ListenAndServe(":1234",nil)
	fmt.Println("server is start : http://localhost:1234")
}
