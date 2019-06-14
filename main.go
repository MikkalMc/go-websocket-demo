package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		//listen and write
		go readWriteRoutine(conn)

		//fake a data stream
		go fakeStream(conn)

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	fmt.Println("running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func readWriteRoutine(conn *websocket.Conn) {
	// Read message from browser
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func fakeStream(conn *websocket.Conn) {
	for {
		time.Sleep(5000 * time.Millisecond)
		msg, err := json.Marshal(rand.Int())
		if err = conn.WriteMessage(1, msg); err != nil {
			return
		}
	}
}
