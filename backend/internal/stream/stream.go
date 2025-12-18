package stream

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type StreamChunk struct {
    Size int
    Data []byte
    Sequence int
}

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
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

func AudioStream(w http.ResponseWriter, r *http.Request) {
    log.Printf("[%v] Opening stream...", r.RemoteAddr)
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("[%v] Failed to start stream: %v", r.RemoteAddr)
    }
    defer conn.Close()
    log.Printf("[%v] Stream started.", conn.RemoteAddr())
    // We should check out PreparedMessage from the websocket library

    // step 1: chunk audio
    // step 2: send audio via websocket
    file, err := os.Open("bleachers.mp3")
    if err != nil {
        log.Printf("[%v] Failed to open file for streaming: %v", conn.RemoteAddr(), err) 
            conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to open file."))
            time.Sleep(10 * time.Second)
            conn.Close()
            return
    }
    defer file.Close()
    streamBuffer := make([]byte, 8192*5)
    sequence := 0
    for {
        chunkSize, err := file.Read(streamBuffer)
        if err == io.EOF {
            log.Printf("[%v] Stream complete!", conn.RemoteAddr())
            conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Stream complete."))
            time.Sleep(10 * time.Second)
            conn.Close()
            return
        }
        if err != nil {
            log.Printf("[%v] Failed to chunk stream: %v", conn.RemoteAddr(), err)
            conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to chunk stream."))
            time.Sleep(10 * time.Second)
            conn.Close()
            return
        }
        chunk := StreamChunk{
            Size: chunkSize,
            Data: streamBuffer[:chunkSize],
            Sequence: sequence,
        }
        message, err := json.Marshal(&chunk)
        if err != nil {
            log.Printf("[%v] Failed to marshal chunk: %v", conn.RemoteAddr(), err)
            conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Failed to chunk stream."))
            time.Sleep(10 * time.Second)
            conn.Close()
            return
        }
        conn.WriteMessage(websocket.BinaryMessage, message)
        sequence++
    }
}
