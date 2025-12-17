package stream

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1204,
    CheckOrigin: checkOrigin,
}

func checkOrigin(r *http.Request) bool {
    return true;
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to open websocket connection: %v", err)
        return
    }
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Failed to read incoming websocket message: %v", err)
            break
        }
        log.Printf("Received: %s", message)
        conn.WriteMessage(messageType, message)
        if err != nil {
            log.Printf("Failed to write a message to websocket: %v", err)
            break
        }
    }
}
